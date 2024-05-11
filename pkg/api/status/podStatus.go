package status

// PodStatus 添加参数请参考 https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.25/#pod-v1-core
type PodStatus struct {
	ContainerStatuses []ContainerStatus
	HostIP            string
	Phase             string
	PodIP             string
}
