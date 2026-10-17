package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

func main() {
	k8sURL := flag.String("k8s", "http://localhost:8001", "Kubernetes API endpoint")
	promURL := flag.String("prom", "http://localhost:9090", "Prometheus API endpoint")
	namespace := flag.String("ns", "default", "Kubernetes namespace")
	port := flag.String("port", "8081", "Cache server port")
	flag.Parse()

	_, _, _ = promURL, namespace, port

	client := createHTTPClient()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	info, err := getNodeInfo(ctx, client, *k8sURL)
	if err != nil {
		return
	}

	log.Printf("Info: %v", info)
}

func getNodeInfo(ctx context.Context, client *http.Client, url string) ([]Node, error) {
	reqURL := fmt.Sprintf("%s/api/v1/nodes", url)

	body, err := doRequest(ctx, client, reqURL, "")
	if err != nil {
		return nil, err
	}

	var res NodeInfo

	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}

	var nodes []Node

	for _, item := range res.Items {
		var runtimes []string
		for _, runtime := range item.NodeStatus.NodeRuntimeHandlers {
			runtimes = append(runtimes, runtime.Name)
		}

		node := Node{
			name:     item.Metadata.Name,
			runtimes: runtimes,
		}

		nodes = append(nodes, node)
	}

	return nodes, nil
}

func fetchPromMetric(ctx context.Context, client *http.Client, promURL, query string) (map[string]string, error) {
	metrics := make(map[string]string)
	reqURL := fmt.Sprintf("%s/api/v1/query?query=%s", promURL, url.QueryEscape(query))

	body, err := doRequest(ctx, client, reqURL, "")
	if err != nil {
		return nil, err
	}

	var res PrometheusResponse
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}

	for _, result := range res.Data.Result {
		pod := result.Metric["pod"]
		if len(result.Value) >= 2 {
			if val, ok := result.Value[1].(string); ok {
				metrics[pod] = val
			}
		}
	}
	return metrics, nil
}

func createHTTPClient() *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return &http.Client{Timeout: 3 * time.Second, Transport: tr}
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
