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

// ListOptions is the query options to a standard REST list call.
type ListOptions struct {
	Kind string `json:"kind,omitempty" protobuf:"bytes,1,opt,name=kind"`

	APIVersion string `json:"apiVersion,omitempty" protobuf:"bytes,2,opt,name=apiVersion"`

	// A selector to restrict the list of returned objects by their labels.
	// Defaults to everything.
	// +optional
	LabelSelector string `json:"labelSelector,omitempty" protobuf:"bytes,1,opt,name=labelSelector"`
	// A selector to restrict the list of returned objects by their fields.
	// Defaults to everything.
	// +optional
	FieldSelector string `json:"fieldSelector,omitempty" protobuf:"bytes,2,opt,name=fieldSelector"`

	// +k8s:deprecated=includeUninitialized,protobuf=6

	// Watch for changes to the described resources and return them as a stream of
	// add, update, and remove notifications. Specify resourceVersion.
	// +optional
	Watch bool `json:"watch,omitempty" protobuf:"varint,3,opt,name=watch"`

	// resourceVersion sets a constraint on what resource versions a request may be served from.
	// See https://kubernetes.io/docs/reference/using-api/api-concepts/#resource-versions for
	// details.
	//
	// Defaults to unset
	// +optional
	ResourceVersion string `json:"resourceVersion,omitempty" protobuf:"bytes,4,opt,name=resourceVersion"`

	// Timeout for the list/watch call.
	// This limits the duration of the call, regardless of any activity or inactivity.
	// +optional
	TimeoutSeconds *int64 `json:"timeoutSeconds,omitempty" protobuf:"varint,5,opt,name=timeoutSeconds"`
}
