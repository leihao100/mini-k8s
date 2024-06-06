package config

import (
	"MiniK8S/pkg/api/meta"
	"MiniK8S/pkg/api/status"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/google/uuid"
)

type Pod struct {
	ApiVersion string           `json:"apiVersion,omitempty"`
	Kind       string           `json:"kind,omitempty"`
	Metadata   meta.ObjectMeta  `json:"metadata,omitempty"`
	Spec       PodSpec          `json:"spec,omitempty"`
	Status     status.PodStatus `json:"status,omitempty"`
}

const PodRestartTimes = 5

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

type PodList struct {
	ApiVersion      string `json:"apiVersion,omitempty"`
	Kind            string `json:"kind,omitempty"`
	ResourceVersion string `json:"resourceVersion,omitempty"`
	Continue        string `json:"continue,omitempty"`
	Items           []Pod  `json:"items"`
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
func (p *Pod) GetName() string {
	return p.Metadata.Name
}

func (p *Pod) SetResourceVersion(version int64) {
	p.Metadata.ResourceVersion = strconv.FormatInt(version, 10)
}
func (p *Pod) GetResourceVersion() int64 {
	res, err := strconv.ParseInt(p.Metadata.ResourceVersion, 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return 0
	}
	return res
}
func (p *Pod) JsonUnmarshalStatus(data []byte) error {
	return json.Unmarshal(data, &(p.Status))
}

func (p *Pod) JsonMarshalStatus() ([]byte, error) {
	return json.Marshal(p.Status)
}
func (p *Pod) SetStatus(s ApiObjectStatus) bool {
	status, ok := s.(*status.PodStatus)
	if ok {
		p.Status = *status
	}
	return ok
}
func (p *Pod) GetStatus() ApiObjectStatus {
	return &p.Status
}
func (p *Pod) Info() {
	fmt.Printf("%-10s\t%-10s\t%-10s\t%-20s\t%-20s\n", "NAME", "UID", "NODE", "STATUS", "IP")
	fmt.Printf("%-10s\t%-10s\t%-10s\t%-20s\t%-20s\n", p.Metadata.Name, p.Metadata.Uid, p.Spec.NodeName, p.Status.Phase, p.Status.PodIP)
}

func (p *PodList) JsonUnmarshal(data []byte) error {
	return json.Unmarshal(data, &p)
}

func (p *PodList) JsonMarshal() ([]byte, error) {
	return json.Marshal(p)
}
func (p *PodList) AppendItems(objects []string) error {
	for _, object := range objects {
		ApiObject := &Pod{}
		err := ApiObject.JsonUnmarshal([]byte(object))
		if err != nil {
			return err
		}
		p.Items = append(p.Items, *ApiObject)
	}
	return nil
}
func (p *PodList) GetItems() []ApiObject {
	var items []ApiObject
	items = make([]ApiObject, 0)
	for _, item := range p.Items {
		items = append(items, &item)
	}
	return items
}
func (p *PodList) Info() {
	fmt.Printf("%-10s\t%-10s\t%10s\t%-20s\t%-20s\n", "NAME", "UID", "NODE", "STATUS", "IP")
	for _, item := range p.Items {
		fmt.Printf("%-10s\t%-10s\t%-10s\t%-20s\t%-20s\n", item.Metadata.Name, item.Metadata.Uid, item.Spec.NodeName, item.Status.Phase, item.Status.PodIP)
	}
}
