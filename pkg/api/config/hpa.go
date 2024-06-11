package config

import (
	"MiniK8S/pkg/api/meta"
	"MiniK8S/pkg/api/status"
	"MiniK8S/pkg/api/types"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/google/uuid"
)

type HorizontalPodAutoscaler struct {
	ApiVersion string                               `json:"apiversion,omitempty"`
	Kind       string                               `json:"kind,omitempty"`
	Metadata   meta.ObjectMeta                      `json:"metadata,omitempty"`
	Spec       HorizontalPodAutoscalerSpec          `json:"spec,omitempty"`
	Status     status.HorizontalPodAutoscalerStatus `json:"status,omitempty"`
}

type HorizontalPodAutoscalerSpec struct {
	MaxReplicas    int32                           `json:"maxReplicas,omitempty"`
	ScaleTargetRef CrossVersionObjectReference     `json:"scaleTargetRef,omitempty"`
	MinReplicas    int32                           `json:"minReplicas,omitempty"`
	Behavior       HorizontalPodAutoscalerBehavior `json:"behavior,omitempty"`
	Metrics        []MetricSpec                    `json:"metrics,omitempty"`
}

/*
 maxReplicas (int32)，必需
maxReplicas 是自动扩缩器可以扩容的副本数的上限。不能小于 minReplicas。
scaleTargetRef (CrossVersionObjectReference)，必需
scaleTargetRef 指向要扩缩的目标资源，用于收集 Pod 的相关指标信息以及实际更改的副本数。
CrossVersionObjectReference 包含足够的信息来让你识别出所引用的资源。
minReplicas (int32)
minReplicas 是自动扩缩器可以缩减的副本数的下限。它默认为 1 个 Pod。 如果启用了 Alpha 特性门控 HPAScaleToZero 并且配置了至少一个 Object 或 External 度量指标， 则 minReplicas 允许为 0。只要至少有一个度量值可用，扩缩就处于活动状态。
behavior (HorizontalPodAutoscalerBehavior)
behavior 配置目标在扩容（Up）和缩容（Down）两个方向的扩缩行为（分别用 scaleUp 和 scaleDown 字段）。 如果未设置，则会使用默认的 HPAScalingRules 进行扩缩容。
HorizontalPodAutoscalerBehavior 配置目标在扩容（Up）和缩容（Down）两个方向的扩缩行为 （分别用 scaleUp 和 scaleDown 字段）。
metrics ([]MetricSpec)
原子性：将在合并时被替换
metrics 包含用于计算预期副本数的规约（将使用所有指标的最大副本数）。 预期副本数是通过将目标值与当前值之间的比率乘以当前 Pod 数来计算的。 因此，使用的指标必须随着 Pod 数量的增加而减少，反之亦然。 有关每种类别的指标必须如何响应的更多信息，请参阅各个指标源类别。 如果未设置，默认指标将设置为 80% 的平均 CPU 利用率。
MetricSpec 指定如何基于单个指标进行扩缩容（一次只能设置 type 和一个其他匹配字段）
*/

type CrossVersionObjectReference struct {
	Kind       string `json:"kind,omitempty"`
	Name       string `json:"name,omitempty"`
	ApiVersion string `json:"apiVersion,omitempty"`
}

/*
scaleTargetRef.kind (string)，必需
kind 是被引用对象的类别；更多信息： https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
scaleTargetRef.name (string)，必需
name 是被引用对象的名称；更多信息：https://kubernetes.io/zh-cn/docs/concepts/overview/working-with-objects/names/#names
scaleTargetRef.apiVersion (string)
apiVersion 是被引用对象的 API 版本。
*/

type HorizontalPodAutoscalerBehavior struct {
	ScaleDown HPAScalingRules `json:"scaleDown,omitempty"`
	ScaleUp   HPAScalingRules `json:"scaleUp,omitempty"`
}

/*
behavior.scaleDown (HPAScalingRules)
scaleDown 是缩容策略。如果未设置，则默认值允许缩减到 minReplicas 数量的 Pod， 具有 300 秒的稳定窗口（使用最近 300 秒的最高推荐值）。
HPAScalingRules 为一个方向配置扩缩行为。在根据 HPA 的指标计算 desiredReplicas 后应用这些规则。 可以通过指定扩缩策略来限制扩缩速度。可以通过指定稳定窗口来防止抖动， 因此不会立即设置副本数，而是选择稳定窗口中最安全的值。
behavior.scaleUp (HPAScalingRules)
scaleUp 是用于扩容的扩缩策略。如果未设置，则默认值为以下值中的较高者：
每 60 秒增加不超过 4 个 Pod
每 60 秒 Pod 数量翻倍
不使用稳定窗口。
*/

type SelectPolicy string

const (
	Max SelectPolicy = "max"
	Min SelectPolicy = "min"
)

type HPAScalingRules struct {
	Policies                   []HPAScalingPolicy `json:"policies,omitempty"`
	SelectPolicy               SelectPolicy       `json:"selectPolicy,omitempty"`
	StabilizationWindowSeconds int32              `json:"stabilizationWindowSeconds,omitempty"`
}

/*
behavior.scaleDown.policies ([]HPAScalingPolicy)
原子性：将在合并时被替换
policies 是可在扩缩容过程中使用的潜在扩缩策略的列表。必须至少指定一个策略，否则 HPAScalingRules 将被视为无效而丢弃。
HPAScalingPolicy 是一个单一的策略，它必须在指定的过去时间间隔内保持为 true。
behavior.scaleDown.selectPolicy (string)
selectPolicy 用于指定应该使用哪个策略。如果未设置，则使用默认值 Max。
behavior.scaleDown.stabilizationWindowSeconds (int32)
stabilizationWindowSeconds 是在扩缩容时应考虑的之前建议的秒数。stabilizationWindowSeconds 必须大于或等于零且小于或等于 3600（一小时）。如果未设置，则使用默认值：
扩容：0（不设置稳定窗口）。
缩容：300（即稳定窗口为 300 秒）。
*/

type PolicyType string

const (
	PolicyPod     PolicyType = "pod"
	PolicyPercent PolicyType = "percent"
)

type HPAScalingPolicy struct {
	Type          PolicyType `json:"type,omitempty"`
	Value         int32      `json:"value,omitempty"`
	PeriodSeconds int32      `json:"periodSeconds,omitempty"`
}

/*
behavior.scaleDown.policies.type (string)，必需
type 用于指定扩缩策略。
behavior.scaleDown.policies.value (int32)，必需
value 包含策略允许的更改量。它必须大于零。
behavior.scaleDown.policies.periodSeconds (int32)，必需
periodSeconds 表示策略应该保持为 true 的时间窗口长度。 periodSeconds 必须大于零且小于或等于 1800（30 分钟）。
*/

type MetricSpec struct {
	Type              string                        `json:"type,omitempty"`
	ContainerResource ContainerResourceMetricSource `json:"containerResource,omitempty"`
	Target            uint64                        `json:"target,omitempty"`
}

/*
metrics.type (string)，必需
type 是指标源的类别。它取值是 “ContainerResource”、“External”、“Object”、“Pods” 或 “Resource” 之一， 每个类别映射到对象中的一个对应的字段。注意：“ContainerResource” 类别在特性门控 HPAContainerMetrics 启用时可用。
metrics.containerResource (ContainerResourceMetricSource)
containerResource 是指 Kubernetes 已知的资源指标（例如在请求和限制中指定的那些）， 描述当前扩缩目标中每个 Pod 中的单个容器（例如 CPU 或内存）。 此类指标内置于 Kubernetes 中，在使用 “pods” 源的、按 Pod 计算的普通指标之外，还具有一些特殊的扩缩选项。 这是一个 Alpha 特性，可以通过 HPAContainerMetrics 特性标志启用。
ContainerResourceMetricSource 指示如何根据请求和限制中指定的 Kubernetes 已知的资源指标进行扩缩容， 此结构描述当前扩缩目标中的每个 Pod（例如 CPU 或内存）。在与目标值比较之前，这些值先计算平均值。 此类指标内置于 Kubernetes 中，并且在使用 “Pods” 源的、按 Pod 统计的普通指标之外支持一些特殊的扩缩选项。 只应设置一种 “target” 类别。
*/
type ContainerResourceMetricSource struct {
	Container string       `json:"container,omitempty"`
	Target    MetricTarget `json:"target,omitempty"`
}

/*
metrics.containerResource.container (string)，必需
container 是扩缩目标的 Pod 中容器的名称。
metrics.containerResource.target (MetricTarget)，必需
target 指定给定指标的目标值。
MetricTarget 定义特定指标的目标值、平均值或平均利用率s
*/
type MetricTarget struct {
	Type               string         `json:"type,omitempty"`
	AverageUtilization int32          `json:"averageUtilization,omitempty"`
	AverageValue       types.Quantity `json:"averageValue,omitempty"`
	Value              types.Quantity `json:"value,omitempty"`
}

/*
metrics.containerResource.target.type (string)，必需
type 表示指标类别是 Utilization、Value 或 AverageValue。
metrics.containerResource.target.averageUtilization (int32)
averageUtilization 是跨所有相关 Pod 的资源指标均值的目标值， 表示为 Pod 资源请求值的百分比。目前仅对 “Resource” 指标源类别有效。
metrics.containerResource.target.averageValue (Quantity)
是跨所有相关 Pod 的指标均值的目标值（以数量形式给出）。
metrics.containerResource.target.value (Quantity)
value 是指标的目标值（以数量形式给出）。
*/

type HorizontalPodAutoscalerList struct {
	ApiVersion      string                    `json:"apiVersion,omitempty"`
	Kind            string                    `json:"kind,omitempty"`
	ResourceVersion string                    `json:"resourceVersion,omitempty"`
	Continue        string                    `json:"continue,omitempty"`
	Items           []HorizontalPodAutoscaler `json:"items"`
}

func (h *HorizontalPodAutoscaler) JsonMarshal() ([]byte, error) {
	return json.Marshal(h)
}

func (h *HorizontalPodAutoscaler) JsonUnmarshal(data []byte) error {
	return json.Unmarshal(data, &h)
}

func (h *HorizontalPodAutoscaler) SetUID(uid uuid.UUID) {
	h.Metadata.Uid = uid
}

func (h *HorizontalPodAutoscaler) GetUID() uuid.UUID {
	return h.Metadata.Uid
}
func (h *HorizontalPodAutoscaler) GetName() string {
	return h.Metadata.Name
}

func (h *HorizontalPodAutoscaler) SetResourceVersion(version int64) {
	h.Metadata.ResourceVersion = strconv.FormatInt(version, 10)
}
func (h *HorizontalPodAutoscaler) GetResourceVersion() int64 {
	res, err := strconv.ParseInt(h.Metadata.ResourceVersion, 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return 0
	}
	return res
}
func (h *HorizontalPodAutoscaler) JsonUnmarshalStatus(data []byte) error {
	return json.Unmarshal(data, &(h.Status))
}

func (h *HorizontalPodAutoscaler) JsonMarshalStatus() ([]byte, error) {
	return json.Marshal(h.Status)
}
func (h *HorizontalPodAutoscaler) SetStatus(s ApiObjectStatus) bool {
	status, ok := s.(*status.HorizontalPodAutoscalerStatus)
	if ok {
		h.Status = *status
	}
	return ok
}
func (h *HorizontalPodAutoscaler) GetStatus() ApiObjectStatus {
	return &h.Status
}
func (h *HorizontalPodAutoscaler) Info() {
	fmt.Printf("%-10s\t%-40s\t%-20s\t%-20s\t%-20s\t%-20s\n", "NAME", "UID", "REFERENCE", "MINPODS", "MAXPODS", "REPLICAS")
	fmt.Printf("%-10s\t%-40s\t%-20s\t%-20d\t%-20d\t%-20d\n", h.Metadata.Name, h.Metadata.Uid, h.Spec.ScaleTargetRef.Kind+"/"+h.Spec.ScaleTargetRef.Name, h.Spec.MinReplicas, h.Spec.MaxReplicas, h.Status.CurrentReplicas)
}
func (h *HorizontalPodAutoscalerList) JsonUnmarshal(data []byte) error {
	return json.Unmarshal(data, &h)
}

func (h *HorizontalPodAutoscalerList) JsonMarshal() ([]byte, error) {
	return json.Marshal(h)
}
func (h *HorizontalPodAutoscalerList) AppendItems(objects []string) error {
	for _, object := range objects {
		ApiObject := &HorizontalPodAutoscaler{}
		err := ApiObject.JsonUnmarshal([]byte(object))
		if err != nil {
			return err
		}
		h.Items = append(h.Items, *ApiObject)
	}
	return nil
}
func (h *HorizontalPodAutoscalerList) GetItems() []ApiObject {
	var items []ApiObject
	items = make([]ApiObject, 0)
	for _, item := range h.Items {
		items = append(items, &item)
	}
	return items
}
func (h *HorizontalPodAutoscalerList) Info() {
	fmt.Printf("%-10s\t%-40s\t%-20s\t%-20s\t%-20s\t%-20s\n", "NAME", "UID", "REFERENCE", "MINPODS", "MAXPODS", "REPLICAS")
	for _, item := range h.Items {
		fmt.Printf("%-10s\t%-40s\t%-20s\t%-20d\t%-20d\t%-20d\n", item.Metadata.Name, item.Metadata.Uid, item.Spec.ScaleTargetRef.Kind+"/"+item.Spec.ScaleTargetRef.Name, item.Spec.MinReplicas, item.Spec.MaxReplicas, item.Status.CurrentReplicas)
	}
}
