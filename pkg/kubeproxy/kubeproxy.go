package kubeproxy

import (
	"MiniK8S/pkg/api/config"
	"MiniK8S/pkg/kubelet"
	"MiniK8S/pkg/kubeproxy/ipInterface"
	ipvsManager "MiniK8S/pkg/kubeproxy/ipvs"
	"github.com/google/uuid"
)

type KubeProxy struct {
	kl            *kubelet.Kubelet
	services      map[uuid.UUID]*config.Service
	ipManager     ipInterface.IP
	serviceToPods map[uuid.UUID][]*config.Pod
}

func NewKubeProxy(kl *kubelet.Kubelet) *KubeProxy {
	return &KubeProxy{
		kl:        kl,
		services:  make(map[uuid.UUID]*config.Service),
		ipManager: ipvsManager.GetIPVS(),
	}
}

func (kp *KubeProxy) CreateService(service *config.Service) {
	kp.services[service.Metadata.Uid] = service
	kp.ipManager.AddService(service)
	pods := kp.SelectPod(service)
	kp.serviceToPods[service.Metadata.Uid] = pods
	for _, pod := range pods {
		kp.ipManager.AddPodToService(service, pod)
	}
}

func (kp *KubeProxy) SelectPod(service *config.Service) []*config.Pod {
	pods := kp.kl.GetPods()
	var targetPods []*config.Pod
	for _, pod := range pods {
		for _, container := range pod.Spec.Containers {
			for s, s2 := range service.Spec.Selector {
				if container.Labels[s] == s2 {
					targetPods = append(targetPods, pod)
				}
			}
		}
	}
	return targetPods
}

func (kp *KubeProxy) RemoveService(service *config.Service) {
	kp.services[service.Metadata.Uid] = nil
	kp.ipManager.RemoveService(service)
	pods := kp.serviceToPods[service.Metadata.Uid]
	for _, pod := range pods {
		kp.ipManager.RemovePodFromService(service, pod)
	}
}

func (kp *KubeProxy) RemovePod(pod *config.Pod) {
	for _, container := range pod.Spec.Containers {
		for _, s2 := range kp.services {
			for k, v := range s2.Spec.Selector {
				if container.Labels[k] == v {
					kp.ipManager.RemovePodFromService(s2, pod)
				}
			}
		}
	}
}

func (kp *KubeProxy) AddPod(pod *config.Pod) {
	for _, container := range pod.Spec.Containers {
		for _, s2 := range kp.services {
			for k, v := range s2.Spec.Selector {
				if container.Labels[k] == v {
					kp.ipManager.AddPodToService(s2, pod)
				}
			}
		}
	}
}

func (kp *KubeProxy) GetSvc() {

}
