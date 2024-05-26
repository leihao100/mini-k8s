package config

import (
	"MiniK8S/pkg/api/status"
	"MiniK8S/pkg/api/types"
	"fmt"

	"github.com/google/uuid"
)

type ApiObject interface {
	JsonUnmarshal([]byte) error
	JsonMarshal() ([]byte, error)
	SetUID(uuid.UUID)
	GetUID() uuid.UUID
	SetResourceVersion(int64)
	GetResourceVersion() int64
	JsonUnmarshalStatus([]byte) error
	JsonMarshalStatus() ([]byte, error)
	SetStatus(ApiObjectStatus) bool
	GetStatus() ApiObjectStatus
	Info()
}
type ApiObjectSpec interface {
}
type ApiObjectStatus interface {
	JsonUnmarshal([]byte) error
	JsonMarshal() ([]byte, error)
}
type ApiObjectList interface {
	JsonUnmarshal([]byte) error
	JsonMarshal() ([]byte, error)
	AppendItems(objects []string) error
	GetItems() any
	Info()
}

// type ErrorApiObject struct {
// 	ApiVersion string          `json:"apiVersion,omitempty"`
// 	Kind       string          `json:"kind,omitempty"`
// 	Metadata   meta.ObjectMeta `json:"metadata,omitempty"`
// 	Spec       ErrorSpec       `json:"spec,omitempty"`
// 	Status     ErrorStatus     `json:"status,omitempty"`
// }

// type ErrorSpec struct {
// }
// type ErrorStatus struct {
// }

// func (e *ErrorApiObject) JsonMarshal() ([]byte, error) {
// 	return json.Marshal(e)
// }

// func (e *ErrorApiObject) JsonUnmarshal(data []byte) error {
// 	return json.Unmarshal(data, &e)
// }

// func (e *ErrorApiObject) SetUID(uid uuid.UUID) {
// 	e.Metadata.Uid = uid
// }

// func (e *ErrorApiObject) GetUID() uuid.UUID {
// 	return e.Metadata.Uid
// }

func NewApiObject(ty types.ApiObjectType) ApiObject {
	switch ty {
	case types.PodObjectType:
		return &Pod{}
	case types.ServiceObjectType:
		return &Service{}
	case types.DeploymentObjectType:
		return &Deployment{}
	case types.HorizontalPodAutoscalerObjectType:
		return &HorizontalPodAutoscaler{}
	case types.NodeObjectType:
		return &Node{}
	}
	panic(fmt.Sprintf("Error ApiObjectType %v", ty))
}

func NewApiObjectStatus(ty types.ApiObjectType) ApiObjectStatus {
	switch ty {
	case types.PodObjectType:
		return &status.PodStatus{}
	case types.ServiceObjectType:
		return &status.ServiceStatus{}
	case types.DeploymentObjectType:
		return &status.DeploymentStatus{}
	case types.HorizontalPodAutoscalerObjectType:
		return &status.HorizontalPodAutoscalerStatus{}
	case types.NodeObjectType:
		return &status.NodeStatus{}
	}
	panic(fmt.Sprintf("Error ApiObjectType %v", ty))
}

func NewApiObjectList(ty types.ApiObjectType) ApiObjectList {
	switch ty {
	case types.PodObjectType:
		return &PodList{}
	case types.ServiceObjectType:
		return &ServiceList{}
	case types.DeploymentObjectType:
		return &DeploymentList{}
	case types.HorizontalPodAutoscalerObjectType:
		return &HorizontalPodAutoscalerList{}
	case types.NodeObjectType:
		return &NodeList{}
	}
	panic(fmt.Sprintf("Error ApiObjectType %v", ty))
}
