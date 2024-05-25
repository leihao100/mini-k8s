package pod

import (
	"MiniK8S/pkg/api/config"
	"MiniK8S/pkg/apiClient"
	"github.com/google/uuid"
)

type PodController struct {
	podClient apiClient.Client
	podMap    map[uuid.UUID]*config.Pod
}

func (pc *PodController) GetPodList() {

}

func (pc *PodController) CreatePod(pod *config.Pod) {
	//pc.podClient.BuildURL(apiClient.Create)
	//buf := pod.Marshal()
	//pod
}
