package config

import (
	"MiniK8S/pkg/api/meta"
	"MiniK8S/pkg/api/selector"
	"MiniK8S/pkg/api/status"

	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"strconv"
)

type Deployment struct {
	ApiVersion string                  `json:"apiVersion,omitempty"`
	Kind       string                  `json:"kind,omitempty"`
	Metadata   meta.ObjectMeta         `json:"metadata,omitempty"`
	Spec       DeploymentSpec          `json:"spec,omitempty"`
	Status     status.DeploymentStatus `json:"status,omitempty"`
}

type DeploymentSpec struct {
	Selector selector.LabelSelector `json:"selector,omitempty"`
	Template PodTemplateSpec        `json:"template,omitempty"`
	Replicas int32                  `json:"replicas,omitempty"`
}

/*
selector (LabelSelector)，必需
供 Pod 所用的标签选择算符。通过此字段选择现有 ReplicaSet 的 Pod 集合， 被选中的 ReplicaSet 将受到这个 Deployment 的影响。此字段必须与 Pod 模板的标签匹配。
template (PodTemplateSpec)，必需
template 描述将要创建的 Pod。template.spec.restartPolicy 唯一被允许的值是 Always。
replicas (int32)
预期 Pod 的数量。这是一个指针，用于辨别显式零和未指定的值。默认为 1
*/

type DeploymentList struct {
	ApiVersion      string       `json:"apiVersion,omitempty"`
	Kind            string       `json:"kind,omitempty"`
	ResourceVersion string       `json:"resourceVersion,omitempty"`
	Continue        string       `json:"continue,omitempty"`
	Items           []Deployment `json:"items"`
}

func (d *Deployment) JsonMarshal() ([]byte, error) {
	return json.Marshal(d)
}

func (d *Deployment) JsonUnmarshal(data []byte) error {
	return json.Unmarshal(data, &d)
}

func (d *Deployment) SetUID(uid uuid.UUID) {
	d.Metadata.Uid = uid
}

func (d *Deployment) GetUID() uuid.UUID {
	return d.Metadata.Uid
}

func (d *Deployment) SetResourceVersion(version int64) {
	d.Metadata.ResourceVersion = strconv.FormatInt(version, 10)
}
func (d *Deployment) GetResourceVersion() int64 {
	res, err := strconv.ParseInt(d.Metadata.ResourceVersion, 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return 0
	}
	return res
}
func (d *Deployment) JsonUnmarshalStatus(data []byte) error {
	return json.Unmarshal(data, &(d.Status))
}

func (d *Deployment) JsonMarshalStatus() ([]byte, error) {
	return json.Marshal(d.Status)
}
func (d *Deployment) SetStatus(s ApiObjectStatus) bool {
	status, ok := s.(*status.DeploymentStatus)
	if ok {
		d.Status = *status
	}
	return ok
}
func (d *Deployment) GetStatus() ApiObjectStatus {
	return &d.Status
}
func (d *Deployment) Info() {
	fmt.Printf("%-10s\t%-10s\t%-10s\t%-20s\n", "NAME", "UID", "DESIRED", "CURRENT")
	fmt.Printf("%-10s\t%-10s\t%-10d\t%-20d\n", d.Metadata.Name, d.Metadata.Uid, d.Spec.Replicas, d.Status.Replicas)
}
func (d *DeploymentList) JsonUnmarshal(data []byte) error {
	return json.Unmarshal(data, &d)
}

func (d *DeploymentList) JsonMarshal() ([]byte, error) {
	return json.Marshal(d)
}
func (d *DeploymentList) AppendItems(objects []string) error {
	for _, object := range objects {
		ApiObject := &Deployment{}
		err := ApiObject.JsonUnmarshal([]byte(object))
		if err != nil {
			return err
		}
		d.Items = append(d.Items, *ApiObject)
	}
	return nil
}
func (d *DeploymentList) GetItems() []ApiObject {
	var items []ApiObject
	items = make([]ApiObject, 0)
	for _, item := range d.Items {
		items = append(items, &item)
	}
	return items
}
func (d *DeploymentList) Info() {
	fmt.Printf("%-10s\t%-10s\t%10s\t%-20s\n", "NAME", "UID", "DESIRED", "CURRENT")
	for _, item := range d.Items {
		fmt.Printf("%-10s\t%-10s\t%-10d\t%-20d\n", item.Metadata.Name, item.Metadata.Uid, item.Spec.Replicas, item.Status.Replicas)
	}
}
