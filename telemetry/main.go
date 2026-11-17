package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func main() {
	k8sURL := flag.String("k8s", "http://localhost:8001", "Kubernetes API endpoint")
	client := &http.Client{Timeout: 0}
	cache := NewCacheStore()
	notifier := NewPodNotifier()
	baseURL := "http://localhost:8001"

	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		initialCtx, cancelInitial := context.WithTimeout(context.Background(), 4*time.Second)
		if err := syncClusterTopology(initialCtx, client, *k8sURL, cache); err != nil {
			log.Printf("Initial sync failed: %v", err)
		}
		cancelInitial()

		for range ticker.C {
			syncCtx, cancelSync := context.WithTimeout(context.Background(), 4*time.Second)

			err := syncClusterTopology(syncCtx, client, *k8sURL, cache)
			if err != nil {
				log.Printf("Cluster sync error: %v", err)
			} else {
				logCacheContents(cache)
			}

			cancelSync()
		}
	}()

	go startPodWatcher(context.Background(), client, baseURL, notifier)

	broker := &SmartBroker{
		cache:    cache,
		notifier: notifier,
		client:   client,
		baseURL:  baseURL,
	}

	log.Printf("Smart Broker proxy listening on http://localhost:8081")
	if err := http.ListenAndServe(":8081", broker); err != nil {
		log.Fatalf("Proxy failed: %v", err)
	}
}

func logCacheContents(cache *CacheStore) {
	topo := cache.Get()

	log.Printf("=== Cluster Topology Cache ===")

	if len(topo.Nodes) == 0 {
		log.Println("Nodes: none discovered")
	} else {
		for nodeName, nodeInfo := range topo.Nodes {
			log.Printf("Node [%s]: %+v", nodeName, nodeInfo)
		}
	}

	if len(topo.CompatibleNodes) == 0 {
		log.Println("Compatible nodes: none")
	} else {
		for runtime, nodes := range topo.CompatibleNodes {
			if runtime == "" {
				continue
			}
			log.Printf("Compatible nodes for [%s] (%d): %v", runtime, len(nodes), nodes)
		}
	}

	if len(topo.WarmPods) == 0 {
		log.Println("Warm pods: none")
	} else {
		for runtime, pods := range topo.WarmPods {
			if runtime == "" {
				continue
			}
			log.Printf("Warm pods for [%s] (%d): %+v", runtime, len(pods), pods)
		}
	}
}

// -- maybe a better idea would be to separate the node from the machine itself so we can look up active nodes based on
// -- runtimes and machines based on runtimes for better performance, then the question becomes, how can I tell a certain
// -- node to boot up a new instance with a specific runtime and then also return the IP of the given pod that was created
func syncClusterTopology(ctx context.Context, client *http.Client, baseURL string, cache *CacheStore) error {
	nodeURL := fmt.Sprintf("%s/api/v1/nodes", baseURL)
	nodeBody, err := doRequest(ctx, client, nodeURL, "")
	if err != nil {
		return err
	}

	var res NodeInfo
	if err := json.Unmarshal(nodeBody, &res); err != nil {
		return err
	}

	podURL := fmt.Sprintf("%s/api/v1/namespaces/default/pods", baseURL)
	podBody, err := doRequest(ctx, client, podURL, "")
	if err != nil {
		return err
	}

	var podInfo PodInfo
	if err := json.Unmarshal(podBody, &podInfo); err != nil {
		return err
	}

	newTopology := buildTopology(res, podInfo)

	cache.Update(newTopology)

	return nil
}

func buildTopology(rawNodes NodeInfo, rawPods PodInfo) Topology {
	nodesMap := make(map[string]Node)
	warmPodsMap := make(map[string][]ActivePod)
	compatibleNodesMap := make(map[string][]string)

	for _, item := range rawNodes.Items {
		nodeName := item.Metadata.Name
		var runtimes []string

		for _, runtime := range item.NodeStatus.NodeRuntimeHandlers {
			runtimes = append(runtimes, runtime.Name)
			compatibleNodesMap[runtime.Name] = append(compatibleNodesMap[runtime.Name], nodeName)
		}

		nodesMap[nodeName] = Node{
			Name:     nodeName,
			Runtimes: runtimes,
			Pods:     []ActivePod{},
		}
	}

	for _, item := range rawPods.Items {
		// Filter: Must be Running and have an assigned IP
		if item.Status.Phase != "Running" || item.Status.PodIP == "" {
			continue
		}

		// Filter: Must be fully Ready
		isReady := false
		for _, cond := range item.Status.Conditions {
			if cond.Type == "Ready" && cond.Status == "True" {
				isReady = true
				break
			}
		}
		if !isReady {
			continue
		}

		runtime := item.Spec.Runtime
		if runtime == "" {
			runtime = "runc"
		}

		activePod := ActivePod{
			Name:    item.Metadata.Name,
			IP:      item.Status.PodIP,
			Runtime: runtime,
		}

		nodeName := item.Spec.NodeName
		if node, exists := nodesMap[nodeName]; exists {
			node.Pods = append(node.Pods, activePod)
			nodesMap[nodeName] = node
			warmPodsMap[runtime] = append(warmPodsMap[runtime], activePod)
		}
	}

	return Topology{
		Nodes:           nodesMap,
		WarmPods:        warmPodsMap,
		CompatibleNodes: compatibleNodesMap,
	}
}

func doRequest(ctx context.Context, client *http.Client, reqURL, token string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, err
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}
