package main

type MachineInfo struct {
	PodInfo     []PodInfo
	RuntimeInfo []RuntimeInfo
}

// --- Pod Structs ---

type PodInfo struct {
	Name         string
	IP           string
	Node         string
	RuntimeClass string
	CPU          string
	Memory       string
}

type Pods struct {
	Items []PodItem `json:"items"`
}

type PodItem struct {
	Metadata PodMetadata `json:"metadata"`
	Spec     PodSpec     `json:"spec"`
	Status   PodStatus   `json:"status"`
}

type PodMetadata struct {
	Name string `json:"name"`
}

type PodSpec struct {
	NodeName         string `json:"nodeName"`
	RuntimeClassName string `json:"runtimeClassName"`
}

type PodStatus struct {
	PodIP string `json:"podIP"`
	Phase string `json:"phase"`
}

// --- Runtime Structs ---

type RuntimeInfo struct {
	Name string
}

type Runtimes struct {
	Items []RuntimeItem `json:"items"`
}

type RuntimeItem struct {
	Metadata RuntimeMetadata `json:"metadata"`
}

type RuntimeMetadata struct {
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
