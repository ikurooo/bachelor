package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

type SmartBroker struct {
	cache    *CacheStore
	notifier *PodNotifier
	client   *http.Client
	baseURL  string
}

func (b *SmartBroker) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	topo := b.cache.Get()

	const payloadThreshold = 500 * 1024
	desiredRuntime := "wasmtime"

	if r.ContentLength > payloadThreshold || r.URL.Query().Get("heavy") == "true" {
		desiredRuntime = "runc"
		log.Printf("Payload size (%d bytes) exceeds threshold. Routing to Docker (runc)", r.ContentLength)
	} else {
		log.Printf("Payload size (%d bytes) is lightweight. Routing to Wasm (wasmtime)", r.ContentLength)
	}

	var targetPod *ActivePod
	if warmPods, exists := topo.WarmPods[desiredRuntime]; exists && len(warmPods) > 0 {
		targetPod = &warmPods[0]
	}

	if targetPod == nil {
		log.Printf("No warm instance found for runtime: [%s]. Initiating event-driven cold start...", desiredRuntime)

		eligibleNodes := topo.CompatibleNodes[desiredRuntime]
		if len(eligibleNodes) == 0 {
			http.Error(w, fmt.Sprintf("503 - No cluster nodes support runtime: %s", desiredRuntime), http.StatusServiceUnavailable)
			return
		}

		targetNodeName := eligibleNodes[0]
		log.Printf("Selected compatible node [%s] for cold-starting [%s]", targetNodeName, desiredRuntime)

		podName, err := triggerColdStart(b.client, b.baseURL, "default", targetNodeName, desiredRuntime)
		if err != nil {
			log.Printf("Failed to trigger cold start pod: %v", err)
			http.Error(w, "503 - Failed to trigger cold start", http.StatusServiceUnavailable)
			return
		}

		podReadyChan := b.notifier.Subscribe(podName)
		defer b.notifier.Unsubscribe(podName, podReadyChan)

		log.Printf("Subscribed to watch stream for cold-start pod [%s]. Waiting for ready event...", podName)

		waitCtx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
		defer cancel()

		var podIP string
		select {
		case <-waitCtx.Done():
			http.Error(w, "504 - Gateway Timeout: Cold start took too long", http.StatusGatewayTimeout)
			return
		case podIP = <-podReadyChan:
			log.Printf("Received instant pod IP via Watch stream: [%s]", podIP)
		}

		targetStr := fmt.Sprintf("http://%s:8080", podIP)
		targetURL, err := url.Parse(targetStr)
		if err != nil {
			http.Error(w, "Internal proxy error parsing target IP", http.StatusInternalServerError)
			return
		}

		log.Printf("Routing queued request for [%s] to freshly warmed Pod at: [%s]", desiredRuntime, targetStr)

		proxy := &httputil.ReverseProxy{
			Rewrite: func(preq *httputil.ProxyRequest) {
				preq.SetURL(targetURL)
				preq.SetXForwarded()
			},
			Transport: &coldStartTransport{
				transport: http.DefaultTransport,
			},
		}

		proxy.ServeHTTP(w, r)
		return
	}

	targetStr := fmt.Sprintf("http://%s:8080", targetPod.IP)
	targetURL, err := url.Parse(targetStr)
	if err != nil {
		http.Error(w, "Internal proxy error parsing target IP", http.StatusInternalServerError)
		return
	}

	log.Printf("Routing warm request for [%s] to Pod: %s (%s)", desiredRuntime, targetPod.Name, targetStr)

	proxy := &httputil.ReverseProxy{
		Rewrite: func(preq *httputil.ProxyRequest) {
			preq.SetURL(targetURL)
			preq.SetXForwarded()
		},
	}

	proxy.ServeHTTP(w, r)
}

func triggerColdStart(client *http.Client, baseURL, namespace, nodeName, runtime string) (string, error) {
	podName := fmt.Sprintf("cold-%s-%d", runtime, time.Now().UnixNano())

	image := "wasm-serverless:latest"
	if runtime == "runc" {
		image = "native-app:latest" // Matches your native container manifest image name!
	}

	podSpec := map[string]any{
		"nodeName": nodeName,
		"containers": []map[string]any{
			{
				"name":            "app",
				"image":           image,
				"imagePullPolicy": "Never",
				"ports": []map[string]any{
					{"containerPort": 8080},
				},

				"readinessProbe": map[string]any{
					"httpGet": map[string]any{
						"path": "/hello",
						"port": 8080,
					},
					"initialDelaySeconds": 0,
					"periodSeconds":       1,
					"successThreshold":    1,
				},
			},
		},
	}

	if runtime != "runc" {
		podSpec["runtimeClassName"] = runtime
	}

	podManifest := map[string]any{
		"apiVersion": "v1",
		"kind":       "Pod",
		"metadata": map[string]any{
			"name":      podName,
			"namespace": namespace,
			"labels": map[string]string{
				"app": "benchmark-app",
			},
		},
		"spec": podSpec,
	}

	payload, err := json.Marshal(podManifest)
	if err != nil {
		return "", err
	}

	reqURL := fmt.Sprintf("%s/api/v1/namespaces/%s/pods", baseURL, namespace)

	req, err := http.NewRequest("POST", reqURL, bytes.NewBuffer(payload))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("kubernetes rejected pod creation (status %d): %s", resp.StatusCode, string(bodyBytes))
	}

	return podName, nil
}

type coldStartTransport struct {
	transport http.RoundTripper
}

func (t *coldStartTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	maxRetries := 10
	var resp *http.Response
	var err error

	for i := range maxRetries {
		resp, err = t.transport.RoundTrip(req)
		if err == nil {
			return resp, nil
		}

		if netErr, ok := errors.AsType[*net.OpError](err); ok {
			var scErr *net.AddrError
			if errors.As(netErr.Err, &scErr) || netErr.Op == "dial" {
				log.Printf("Pod is booting up, connection refused. Retrying in 100ms... (Attempt %d/%d)", i+1, maxRetries)
				time.Sleep(100 * time.Millisecond)
				continue
			}
		}

		return nil, err
	}

	return nil, err
}
