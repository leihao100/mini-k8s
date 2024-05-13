package status

import "MiniK8S/pkg/api/types"

type ServiceStatus struct {
	Conditions   []ServiceCondition `json:"conditions,omitempty"`
	LoadBalancer LoadBalancerStatus `json:"loadBalancer,omitempty"`
}

/*
Patch strategy: 在 type 上合并
Map: 键类型的唯一值将在合并期间保留
服务的当前状态。
condition 包含此 API 资源某一方面当前的状态详细信息。
loadBalancer (LoadBalancerStatus)
loadBalancer 包含负载均衡器的当前状态（如果存在）。
LoadBalancerStatus 表示负载均衡器的状态。
*/
type ServiceCondition struct {
	LastTransitionTime types.Time `json:"lastTransitionTime,omitempty"`
	Message            string     `json:"message,omitempty"`
	Reason             string     `json:"reason,omitempty"`
	Status             string     `json:"status,omitempty"`
	Type               string     `json:"type,omitempty"`
	ObservedGeneration int64      `json:"observedGeneration,omitempty"`
}

/*
conditions.lastTransitionTime（Time），必需
lastTransitionTime 是状况最近一次状态转化的时间。 变化应该发生在下层状况发生变化的时候。如果不知道下层状况发生变化的时间， 那么使用 API 字段更改的时间是可以接受的。
Time 是 time.Time 的包装类，支持正确地序列化为 YAML 和 JSON。 为 time 包提供的许多工厂方法提供了包装类。
conditions.message (string)，必需
message 是人类可读的消息，有关转换的详细信息，可以是空字符串。
conditions.reason (string)，必需
reason 包含一个程序标识符，指示 condition 最后一次转换的原因。 特定条件类型的生产者可以定义该字段的预期值和含义，以及这些值是否被视为有保证的 API。 该值应该是 CamelCase 字符串且不能为空。
conditions.status (string)，必需
condition 的状态，True、False、Unknown 之一。
conditions.type (string)，必需
CamelCase 或 foo.example.com/CamelCase 中的条件类型。
conditions.observedGeneration (int64)
observedGeneration 表示设置 condition 基于的 .metadata.generation 的过期次数。 例如，如果 .metadata.generation 当前为 12，但 .status.conditions[x].observedGeneration 为 9， 则 condition 相对于实例的当前状态已过期。
*/

type LoadBalancerStatus struct {
	Ingress []LoadBalancerIngress `json:"ingress,omitempty"`
}

/*
loadBalancer.ingress ([]LoadBalancerIngress)
ingress 是一个包含负载均衡器 Ingress 点的列表。Service 的流量需要被发送到这些 Ingress 点。
*/

type LoadBalancerIngress struct {
	Hostname string       `json:"hostname,omitempty"`
	Ip       string       `json:"ip,omitempty"`
	IpMode   string       `json:"ipMode,omitempty"`
	Ports    []PortStatus `json:"ports,omitempty"`
}

/*
loadBalancer.ingress.hostname (string)
hostname 是为基于 DNS 的负载均衡器 Ingress 点（通常是 AWS 负载均衡器）设置的。
loadBalancer.ingress.ip (string)
ip 是为基于 IP 的负载均衡器 Ingress 点（通常是 GCE 或 OpenStack 负载均衡器）设置的。
loadBalancer.ingress.ipMode (string)
  ipMode 指定负载平衡器 IP 的行为方式，并且只能在设置了 ip 字段时指定。
将其设置为 VIP 表示流量将传送到节点，并将目标设置为负载均衡器的 IP 和端口。 将其设置为 Proxy 表示将流量传送到节点或 Pod，并将目标设置为节点的 IP 和节点端口或 Pod 的 IP 和端口。 服务实现可以使用此信息来调整流量路由。
loadBalancer.ingress.ports ([]PortStatus)
Atomic：将在合并期间被替换
ports 是 Service 的端口列表。如果设置了此字段，Service 中定义的每个端口都应该在此列表中。
*/
