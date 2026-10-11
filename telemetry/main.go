package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"text/tabwriter"
	"time"
)

func main() {
	k8sURL := flag.String("k8s", "http://localhost:8001", "Kubernetes API endpoint")
	promURL := flag.String("prom", "http://localhost:9090", "Prometheus API endpoint")
	namespace := flag.String("ns", "default", "Kubernetes namespace to inspect")
	flag.Parse()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client := createHTTPClient()

	token, err := getServiceAccountToken()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\n==================================================================")
	fmt.Println("             K3s Telemetry & Runtime CLI Inspector               ")
	fmt.Println("==================================================================")

	runtimes, err := fetchRuntimeClasses(ctx, client, *k8sURL, token)
	if err != nil {
		fmt.Printf("RuntimeClass check failed: %v\n", err)
	} else {
		fmt.Printf("Available Cluster Runtimes : %v\n", runtimes)
	}

	pods, err := fetchPods(ctx, client, *k8sURL, token, *namespace)
	if err != nil {
		fmt.Printf("Failed to query Pods from K8s: %v\n", err)
		os.Exit(1)
	}

	cpuMetrics, err := fetchPromMetric(ctx, client, *promURL, `sum(rate(container_cpu_usage_seconds_total{pod=~"native-app|wasm-app"}[1m])) by (pod)`)
	memMetrics, err := fetchPromMetric(ctx, client, *promURL, `sum(container_memory_working_set_bytes{pod=~"native-app|wasm-app"}) by (pod)`)

	fmt.Println("\n--- Pod Telemetry Overview ---")
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "POD NAME\tPOD IP\tNODE\tRUNTIME CLASS\tCPU (CORES)\tMEMORY (BYTES)")
	fmt.Fprintln(w, "--------\t------\t----\t-------------\t-----------\t--------------")

	for _, pod := range pods {
		cpu := cpuMetrics[pod.Name]
		if cpu == "" {
			cpu = "N/A"
		}

		mem := memMetrics[pod.Name]
		if mem == "" {
			mem = "N/A"
		}

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
			pod.Name, pod.IP, pod.Node, pod.RuntimeClass, cpu, mem)
	}
	w.Flush()
	fmt.Println()
}

func createHTTPClient() *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return &http.Client{Timeout: 3 * time.Second, Transport: tr}
}

func getServiceAccountToken() (string, error) {
	tokenBytes, err := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/token")
	if err != nil {
		return "", errors.New("unable to read service account token")
	}
	return string(tokenBytes), nil
}

// still unimplemented
func fetchData(ctx context.Context, client *http.Client, k8sURL string, token string, namespace string, cache *Cache) error {
	newMachineInfoCache := make(map[string]MachineInfo)
	var machineInfo MachineInfo

	runtimes, err := fetchRuntimeClasses(ctx, client, k8sURL, token)
	if err == nil {
		machineInfo.RuntimeInfo = runtimes
	}

	newMachineInfoCache[namespace] = machineInfo
	cache.Update(newMachineInfoCache)
	return nil
}

func fetchRuntimeClasses(ctx context.Context, client *http.Client, k8sURL, token string) ([]RuntimeInfo, error) {
	reqURL := fmt.Sprintf("%s/apis/node.k8s.io/v1/runtimeclasses", k8sURL)
	body, err := doRequest(ctx, client, reqURL, token)
	if err != nil {
		return nil, err
	}

	var res Runtimes
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}

	// Kubernetes does not automatically create RuntimeClass objects for the default container runtime
	runtimes := []RuntimeInfo{
		{"runc (default)"},
	}
	for _, item := range res.Items {
		runtimes = append(runtimes, RuntimeInfo{item.Metadata.Name})
	}
	return runtimes, nil
}

func fetchPods(ctx context.Context, client *http.Client, k8sURL, token, ns string) ([]PodInfo, error) {
	reqURL := fmt.Sprintf("%s/api/v1/namespaces/%s/pods", k8sURL, ns)
	body, err := doRequest(ctx, client, reqURL, token)
	if err != nil {
		return nil, err
	}

	var res Pods
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}

	var pods []PodInfo
	for _, item := range res.Items {
		if item.Status.Phase != "Running" {
			continue
		}
		rc := item.Spec.RuntimeClassName
		if rc == "" {
			rc = "native-runc"
		}

		pods = append(pods, PodInfo{
			Name:         item.Metadata.Name,
			IP:           item.Status.PodIP,
			Node:         item.Spec.NodeName,
			RuntimeClass: rc,
		})
	}
	return pods, nil
}

func fetchPromMetric(ctx context.Context, client *http.Client, promURL, query string) (map[string]string, error) {
	metrics := make(map[string]string)
	reqURL := fmt.Sprintf("%s/api/v1/query?query=%s", promURL, url.QueryEscape(query))

	body, err := doRequest(ctx, client, reqURL, "")
	if err != nil {
		return nil, errors.New("Error querying PrometheusL: " + err.Error())
	}

	var res PrometheusResponse
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, errors.New("Error unmarshalling Prometheus: " + err.Error())
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
