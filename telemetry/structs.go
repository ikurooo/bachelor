package main

type Topology struct {
	Nodes           map[string]Node
	WarmPods        map[string][]ActivePod
	CompatibleNodes map[string][]string
}

type Node struct {
	Name     string      `json:"name"`
	Runtimes []string    `json:"runtimes"`
	Pods     []ActivePod `json:"pods"`
}

type ActivePod struct {
	Name    string `json:"name"`
	IP      string `json:"ip"`
	Runtime string `json:"runtime"`
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

type PodInfo struct {
	Items []PodItems `json:"items"`
}

type PodItems struct {
	Metadata PodMetadata `json:"metadata"`
	Spec     PodSpec     `json:"spec"`
	Status   PodStatus   `json:"status"`
}

type PodMetadata struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

type PodSpec struct {
	NodeName string `json:"nodeName"`
	Runtime  string `json:"runtimeClassName"`
}

type PodStatus struct {
	Phase      string         `json:"phase"`
	PodIP      string         `json:"podIP"`
	Conditions []PodCondition `json:"conditions"`
}

type PodCondition struct {
	Type   string `json:"type"`
	Status string `json:"status"`
}
