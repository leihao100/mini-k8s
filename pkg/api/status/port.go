package status

type PortStatus struct {
	Port     int32  `json:"port,omitempty"`
	Protocol string `json:"protocol,omitempty"`
	Error    string `json:"error,omitempty"`
}

/*
loadBalancer.ingress.ports.port (int32)，必需
port 是所记录的服务端口状态的端口号。
loadBalancer.ingress.ports.protocol (string)，必需
protocol 是所记录的服务端口状态的协议。支持的值为：TCP、UDP、SCTP。
loadBalancer.ingress.ports.error (string)
error 是记录 Service 端口的问题。 错误的格式应符合以下规则：
内置错误原因应在此文件中指定，应使用 CamelCase 名称。
云提供商特定错误原因的名称必须符合格式 foo.example.com/CamelCase。
*/
