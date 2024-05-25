package config

import (
	"MiniK8S/pkg/api/meta"
	"MiniK8S/pkg/api/status"
	"encoding/json"

	"github.com/google/uuid"
)

type Pod struct {
	ApiVersion string           `json:"apiVersion,omitempty"`
	Kind       string           `json:"kind,omitempty"`
	Metadata   meta.ObjectMeta  `json:"metadata,omitempty"`
	Spec       PodSpec          `json:"spec,omitempty"`
	Status     status.PodStatus `json:"status,omitempty"`
}

/*
API文档中描述如下
Field	Description
apiVersion
string	APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources
kind
string	Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
metadata
ObjectMeta	Standard object's metadata. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
spec
PodSpec	Specification of the desired behavior of the pod. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status
status
PodStatus	Most recently observed status of the pod. This data may not be up to date. Populated by the system. Read-only. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status
*/

type PodSpec struct {
	Containers     []Container       `json:"containers,omitempty"`
	InitContainers []Container       `json:"initContainers,omitempty"`
	NodeName       string            `json:"nodeName,omitempty"`
	ExposedPorts   []string          `json:"exposedPorts,omitempty"`
	Volumes        []string          `json:"volumes,omitempty"`
	BindPorts      map[string]string `json:"bindPorts,omitempty"`
}

type PodTemplateSpec struct {
	Metadata meta.ObjectMeta `json:"metadata,omitempty"`
	Spec     PodSpec         `json:"spec,omitempty"`
}

func (p *Pod) JsonMarshal() ([]byte, error) {
	return json.Marshal(p)
}

func (p *Pod) JsonUnmarshal(data []byte) error {
	return json.Unmarshal(data, &p)
}

func (p *Pod) SetUID(uid uuid.UUID) {
	p.Metadata.Uid = uid
}

func (p *Pod) GetUID() uuid.UUID {
	return p.Metadata.Uid
}
