package config

import (
	"MiniK8S/pkg/api/meta"
	"MiniK8S/pkg/api/status"
)

type Service struct {
	ApiVersion string               `json:"apiVersion,omitempty"`
	Kind       string               `json:"kind,omitempty"`
	Metadata   meta.ObjectMeta      `json:"metadata,omitempty"`
	Spec       ServiceSpec          `json:"spec,omitempty"`
	Status     status.ServiceStatus `json:"status,omitempty"`
}

type ServiceSpec struct {
	Selector  map[string]string `json:"selector,omitempty"`
	Ports     []ServicePort     `json:"ports,omitempty"`
	Type      string            `json:"type,omitempty"`
	ClusterIP string            `json:"clusterIP,omitempty"`
}

/*
selector (map[string]string)
将 Service 流量路由到具有与此 selector 匹配的标签键值对的 Pod。 如果为空或不存在，
则假定该服务有一个外部进程管理其端点，Kubernetes 不会修改该端点。 仅适用于 ClusterIP、
NodePort 和 LoadBalancer 类型。如果类型为 ExternalName，则忽略。 更多信息：
https://kubernetes.io/docs/concepts/services-networking/service/

ports ([]ServicePort)
端口映射，详见ServicePort
*/

type ServicePort struct {
	Port       int32  `json:"port,omitempty"`
	TargetPort int32  `json:"targetPort,omitempty"`
	Protocol   string `json:"protocol,omitempty"`
	Name       string `json:"name,omitempty"`
	NodePort   int32  `json:"nodePort,omitempty"`
}

/*
ports.port (int32)，必需
Service 将公开的端口。
ports.targetPort (IntOrString)
在 Service 所针对的 Pod 上要访问的端口号或名称。 编号必须在 1 到 65535 的范围内。名称必须是 IANA_SVC_NAME。 如果此值是一个字符串，将在目标 Pod 的容器端口中作为命名端口进行查找。 如果未指定字段，则使用 port 字段的值（直接映射）。 对于 clusterIP 为 None 的服务，此字段将被忽略， 应忽略不设或设置为 port 字段的取值。 更多信息： https://kubernetes.io/docs/concepts/services-networking/service/#defining-a-service
IntOrString 是一种可以保存 int32 或字符串的类型。 在 JSON 或 YAML 编组和解组中使用时，它会生成或使用内部类型。 例如，这允许您拥有一个可以接受名称或数字的 JSON 字段。
ports.protocol (string)
此端口的 IP 协议。支持 “TCP”、“UDP” 和 “SCTP”。默认为 TCP。
ports.name (string)
Service 中此端口的名称。这必须是 DNS_LABEL。 ServiceSpec 中的所有端口的名称都必须唯一。 在考虑 Service 的端点时，这一字段值必须与 EndpointPort 中的 name 字段相同。 如果此服务上仅定义一个 ServicePort，则为此字段为可选。
ports.nodePort (int32)
当类型为 NodePort 或 LoadBalancer 时，Service 公开在节点上的端口， 通常由系统分配。如果指定了一个在范围内且未使用的值，则将使用该值，否则操作将失败。 如果在创建的 Service 需要该端口时未指定该字段，则会分配端口。 如果在创建不需要该端口的 Service时指定了该字段，则会创建失败。 当更新 Service 时，如果不再需要此字段（例如，将类型从 NodePort 更改为 ClusterIP），这个字段将被擦除。 更多信息： https://kubernetes.io/docs/concepts/services-networking/service/#type-nodeport
*/
