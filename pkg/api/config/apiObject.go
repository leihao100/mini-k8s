package config

import (
	"MiniK8S/pkg/api/meta"
	"MiniK8S/pkg/api/types"
	"encoding/json"

	"github.com/google/uuid"
)

type ApiObject interface {
	JsonUnmarshal(data []byte) error
	JsonMarshal() ([]byte, error)
	SetUID(uid uuid.UUID)
	GetUID() uuid.UUID
}
type ApiObjectSpec interface {
}
type ApiObjectStatus interface {
}

type ErrorApiObject struct {
	ApiVersion string          `json:"apiVersion,omitempty"`
	Kind       string          `json:"kind,omitempty"`
	Metadata   meta.ObjectMeta `json:"metadata,omitempty"`
	Spec       ErrorSpec       `json:"spec,omitempty"`
	Status     ErrorStatus     `json:"status,omitempty"`
}

type ErrorSpec struct {
}
type ErrorStatus struct {
}

func (e *ErrorApiObject) JsonMarshal() ([]byte, error) {
	return json.Marshal(e)
}

func (e *ErrorApiObject) JsonUnmarshal(data []byte) error {
	return json.Unmarshal(data, &e)
}

func (e *ErrorApiObject) SetUID(uid uuid.UUID) {
	e.Metadata.Uid = uid
}

func (e *ErrorApiObject) GetUID() uuid.UUID {
	return e.Metadata.Uid
}

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
	return &ErrorApiObject{}
}
