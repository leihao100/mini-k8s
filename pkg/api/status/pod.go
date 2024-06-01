package status

import "encoding/json"

type PodPhase string

const (
	PodPending   PodPhase = "Pending"
	PodRunning   PodPhase = "Running"
	PodSucceeded PodPhase = "Succeeded"
	PodFailed    PodPhase = "Failed"
	PodUnknown   PodPhase = "Unknown"
)

// PodStatus 添加参数请参考 https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.25/#pod-v1-core
type PodStatus struct {
	ContainerStatuses []ContainerStatus `json:"containerStatuses,omitempty"`
	HostIP            string            `json:"hostIP,omitempty"`
	Phase             string            `json:"phase,omitempty"`
	PodIP             string            `json:"podIP,omitempty"`
}

func (p *PodStatus) JsonMarshal() ([]byte, error) {
	return json.Marshal(p)
}

func (p *PodStatus) JsonUnmarshal(data []byte) error {
	return json.Unmarshal(data, &p)
}
