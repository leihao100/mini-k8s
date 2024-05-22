package ipInterface

import "MiniK8S/pkg/api/config"

// IP 考虑到ipvs不一定可靠，先定义接口方便后续替换
type IP interface {
	AddService(service *config.Service)
	RemoveService(service *config.Service)
	AddPodToService(service *config.Service, pod *config.Pod)
	RemovePodFromService(service *config.Service, pod *config.Pod)
}
