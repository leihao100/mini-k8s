package kubelet

import (
	"MiniK8S/pkg/kubelet/cri"
	"MiniK8S/pkg/kubelet/pod"
	"MiniK8S/pkg/util/config/containerConfig"
	"MiniK8S/pkg/util/config/podConfig"
	"MiniK8S/pkg/util/status"
	"fmt"
)

const pauseName = "mirrorgooglecontainers/pause:latest"

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

func (k *Kubelet) CreatePodPause(pod podConfig.PodConfig) {
	containername := pod.Meta.Namespace + "_" + pod.Meta.Name + "_" + "pause"
	container := containerConfig.ContainerConfig{
		Name:         containername,
		Args:         nil,
		Cmd:          nil,
		Entrypoint:   nil,
		Env:          nil,
		Image:        pauseName,
		Volumes:      nil,
		Labels:       nil,
		PortBindings: nil,
		VolumesFrom:  nil,
		Binds:        nil,
		NetworkMode:  "",
		CPULimit:     0,
		MemLimit:     0,
	}
	response, err := k.cli.CreateContainer(container, containername)
	if err != nil {
		panic(err)
	}
	k.podManager.AddContainer(pod.Meta.Name, response.ID)
}

func (k *Kubelet) MakePod(pod podConfig.PodConfig) {
	podStatus := status.PodStatus{
		ContainerStatuses: nil,
		HostIP:            "",
		Phase:             "",
		PodIP:             "",
	}
	k.CreatePodPause(pod)
	pod.Status = podStatus
	containers := pod.Spec.Containers
	for _, container := range containers {
		containerName := pod.Meta.Namespace + "_" + pod.Meta.Name + "_" + container.Name
		response, err := k.cli.CreateContainer(container, containerName)
		if err != nil {
			panic(err)
			fmt.Println("error:", err)
		}
		k.cli.StartContainer(response.ID)
	}
}

func (k *Kubelet) getPods() []*podConfig.PodConfig {
	return k.podManager.GetPods()
}
