package status

import (
	"MiniK8S/pkg/api/types"
	"encoding/json"
)

type HorizontalPodAutoscalerStatus struct {
	DesiredReplicas int32      `json:"desiredReplicas,omitempty"`
	CurrentReplicas int32      `json:"currentReplicas,omitempty"`
	LastScaleTime   types.Time `json:"lastScaleTime,omitempty"`
}

/*
desiredReplicas (int32)，必需
desiredReplicas 是此自动扩缩器管理的 Pod 的所期望的副本数，由自动扩缩器最后计算。
currentReplicas (int32)
currentReplicas 是此自动扩缩器管理的 Pod 的当前副本数，如自动扩缩器最后一次看到的那样。
lastScaleTime (Time)
lastScaleTime 是 HorizontalPodAutoscaler 上次扩缩 Pod 数量的时间，自动扩缩器使用它来控制更改 Pod 数量的频率。
*/
func (h *HorizontalPodAutoscalerStatus) JsonMarshal() ([]byte, error) {
	return json.Marshal(h)
}

func (h *HorizontalPodAutoscalerStatus) JsonUnmarshal(data []byte) error {
	return json.Unmarshal(data, &h)
}
