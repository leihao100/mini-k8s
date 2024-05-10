package kubelet

import (
	"MiniK8S/pkg/kubelet/cri"
	"MiniK8S/pkg/kubelet/pod"
	"MiniK8S/pkg/util/config/podConfig"
	"fmt"
)

type Kubelet struct {
	cli        cri.Client
	podManager *pod.PodManager
}

func (k *Kubelet) Run() {
	//cli, _ := cri.GetClient()
	var err error
	k.cli, err = cri.GetClient()
	if err != nil {
		panic(err)
		fmt.Println("error:", err)
	}
	k.podManager = pod.NewPodManager()

}
func (k *Kubelet) Stop() {

}

func (k *Kubelet) MakePod(pod podConfig.PodConfig) {
	containers := pod.Spec.Containers
	for _, container := range containers {
		containerName := pod.Meta.Namespace + "_" + pod.Meta.Name + "_" + container.Name
		_, err := k.cli.CreateContainer(container, containerName)
		if err != nil {
			panic(err)
			fmt.Println("error:", err)
		}
		//append(pod.Status.ContainerStatuses, status.ContainerStatus{
		//	Name:         containerName,
		//	ContainerID:  createContainer.ID,
		//	ImageID:      "",
		//	Image:        "",
		//	State:        types.ContainerState{},
		//	Started:      false,
		//	RestartCount: 0,
		//})

	}
}
