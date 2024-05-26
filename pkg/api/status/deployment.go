package status

import (
	"MiniK8S/pkg/api/types"
	"encoding/json"
)

type DeploymentStatus struct {
	Replicas            int32                 `json:"replicas,omitempty"`
	AvailableReplicas   int32                 `json:"availableReplicas,omitempty"`
	ReadyReplicas       int32                 `json:"readyReplicas,omitempty"`
	UnavailableReplicas int32                 `json:"unavailableReplicas,omitempty"`
	UpdatedReplicas     int32                 `json:"updatedReplicas,omitempty"`
	Conditions          []DeploymentCondition `json:"conditions,omitempty"`
	ObservedGeneration  int64                 `json:"observedGeneration,omitempty"`
}

/*
replicas (int32)
此部署所针对的（其标签与选择算符匹配）未终止 Pod 的总数。
availableReplicas (int32)
此部署针对的可用（至少 minReadySeconds 才能就绪）的 Pod 总数。
readyReplicas (int32)
readyReplicas 是此 Deployment 在就绪状况下处理的目标 Pod 数量。
unavailableReplicas (int32)
此部署针对的不可用 Pod 总数。这是 Deployment 具有 100% 可用容量时仍然必需的 Pod 总数。 它们可能是正在运行但还不可用的 Pod，也可能是尚未创建的 Pod。
updatedReplicas (int32)
此 Deployment 所针对的未终止 Pod 的总数，这些 Pod 采用了预期的模板规约。
conditions ([]DeploymentCondition)
补丁策略：按照键 type 合并
表示 Deployment 当前状态的最新可用观测值。
observedGeneration (int64)
Deployment 控制器观测到的代数（Generation）
*/

type DeploymentCondition struct {
	Status             string     `json:"status,omitempty"`
	Type               string     `json:"yype,omitempty"`
	LastTransitionTime types.Time `json:"lastTransitionTime,omitempty"`
	LastUpdateTime     types.Time `json:"lastUpdateTime,omitempty"`
	Message            string     `json:"message,omitempty"`
	Reason             string     `json:"reason,omitempty"`
}

/*
conditions.status (string)，必需
状况的状态，取值为 True、False 或 Unknown 之一。
conditions.type (string)，必需
Deployment 状况的类型。
conditions.lastTransitionTime (Time)
状况上次从一个状态转换为另一个状态的时间。
Time 是对 time.Time 的封装。Time 支持对 YAML 和 JSON 进行正确封包。 为 time 包的许多函数方法提供了封装器。
conditions.lastUpdateTime (Time)
上次更新此状况的时间。
Time 是对 time.Time 的封装。Time 支持对 YAML 和 JSON 进行正确封包。 为 time 包的许多函数方法提供了封装器。
conditions.message (string)
这是一条人类可读的消息，指示有关上次转换的详细信息。
conditions.reason (string)
状况上次转换的原因。
*/

func (d *DeploymentStatus) JsonMarshal() ([]byte, error) {
	return json.Marshal(d)
}

func (d *DeploymentStatus) JsonUnmarshal(data []byte) error {
	return json.Unmarshal(data, &d)
}
