package main

type Node struct {
	name     string
	runtimes []string
}

type NodeInfo struct {
	Items []NodeItems `json:"items"`
}

type NodeItems struct {
	Metadata   NodeMetadata `json:"metadata"`
	NodeStatus NodeStatus   `json:"status"`
}

type NodeMetadata struct {
	Name string `json:"name"`
}

type NodeStatus struct {
	NodeRuntimeHandlers []NodeRuntimeHandlers `json:"runtimeHandlers"`
}

type NodeRuntimeHandlers struct {
	Name string `json:"name"`
}

// --- Prometheus Structs ---

type PrometheusResponse struct {
	Status string         `json:"status"`
	Data   PrometheusData `json:"data"`
}

type PrometheusData struct {
	Result []PrometheusResult `json:"result"`
}

type PrometheusResult struct {
	Metric map[string]string `json:"metric"`
	Value  []any             `json:"value"`
}
