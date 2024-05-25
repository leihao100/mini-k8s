package config

import (
	"MiniK8S/pkg/api/meta"
	"MiniK8S/pkg/api/selector"
	"MiniK8S/pkg/api/status"
	"encoding/json"

	"github.com/google/uuid"
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
