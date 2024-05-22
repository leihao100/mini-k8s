package kubelet

import (
	"MiniK8S/pkg/api/config"
	"MiniK8S/pkg/api/status"
	"MiniK8S/pkg/kubelet/cri"
	"MiniK8S/pkg/kubelet/pod"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/google/uuid"
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

func (k *Kubelet) CreatePodPause(pod *config.Pod) string {
	uid := pod.Metadata.Uid.String()
	name := pod.Metadata.Namespace + "_" + pod.Metadata.Name + "_pause_" + uid
	container := config.Container{
		Name:         name,
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
	//此处需优先创建并启动pause，否则网络将无法配置
	response, err := k.cli.CreatePause(container, name)
	if err != nil {
		panic(err)
	}
	newContainerStatus := status.ContainerStatus{
		Name:         name,
		ContainerID:  response.ID,
		ImageID:      "",
		Image:        container.Image,
		State:        types.ContainerState{},
		Started:      false,
		RestartCount: 0,
	}
	pod.Status.ContainerStatuses = append(pod.Status.ContainerStatuses, newContainerStatus)
	k.podManager.AddContainer(pod.Metadata.Uid, name, response.ID)
	return response.ID
}

func (k *Kubelet) MakePod(pod *config.Pod) {
	pod.Metadata.Uid, _ = uuid.NewUUID()
	k.podManager.AddPod(pod.Metadata.Uid, k.podManager.MakePodName(pod), pod)
	podStatus := status.PodStatus{
		ContainerStatuses: nil,
		HostIP:            "",
		Phase:             "",
		PodIP:             "",
	}
	pauseID := k.CreatePodPause(pod)
	k.cli.StartContainer(pauseID)
	pod.Status = podStatus
	containers := pod.Spec.Containers
	for _, container := range containers {
		containerName := pod.Metadata.Namespace + "_" + pod.Metadata.Name + "_" + container.Name + "_" + pod.Metadata.Uid.String()
		container.Pause = pauseID
		response, err := k.cli.CreateContainer(container, containerName)
		if err != nil {
			panic(err)
			fmt.Println("error:", err)
		}
		k.cli.StartContainer(response.ID)
		newContainerStatus := status.ContainerStatus{
			Name:         containerName,
			ContainerID:  response.ID,
			ImageID:      "",
			Image:        container.Image,
			State:        types.ContainerState{},
			Started:      false,
			RestartCount: 0,
		}
		pod.Status.ContainerStatuses = append(pod.Status.ContainerStatuses, newContainerStatus)
		fmt.Println(len(pod.Status.ContainerStatuses))
	}
}

func (k *Kubelet) GetPods() []*config.Pod {
	return k.podManager.GetPods()
}

/*
	type ContainerJSONBase struct {
	    ID              string `json:"ID"`
	    Created         string
	    Path            string
	    Args            []string
	    State           *ContainerState
	    Image           string
	    ResolvConfPath  string
	    HostnamePath    string
	    HostsPath       string
	    LogPath         string
	    Node            *ContainerNode `json:",omitempty"`
	    Name            string
	    RestartCount    int
	    Driver          string
	    Platform        string
	    MountLabel      string
	    ProcessLabel    string
	    AppArmorProfile string
	    ExecIDs         []string
	    HostConfig      *container.HostConfig
	    GraphDriver     GraphDriverData
	    SizeRw          *int64 `json:",omitempty"`
	    SizeRootFs      *int64 `json:",omitempty"`
	}

	type ContainerState struct {
	    Status     string
	    Running    bool
	    Paused     bool
	    Restarting bool
	    OOMKilled  bool
	    Dead       bool
	    Pid        int
	    ExitCode   int
	    Error      string
	    StartedAt  string
	    FinishedAt string
	    Health     *Health `json:",omitempty"`
	}
*/
func (k *Kubelet) UpdatePodStatusByID(id uuid.UUID) {
	pod := k.podManager.GetPodById(id)
	fmt.Println(len(pod.Status.ContainerStatuses))
	containerStatus := pod.Status.ContainerStatuses
	for i, Status := range containerStatus {
		json, err := k.cli.ContainerStatus(Status.ContainerID)
		if err != nil {
			return
		}
		//fmt.Println("now is"+json.Name+"and its state is", json.State.Running)
		pod.Status.ContainerStatuses[i] = status.ContainerStatus{
			State: types.ContainerState{
				Status:     json.State.Status,
				Running:    json.State.Running,
				Paused:     json.State.Paused,
				Restarting: json.State.Restarting,
				OOMKilled:  json.State.OOMKilled,
				Dead:       json.State.Dead,
				Pid:        json.State.Pid,
				ExitCode:   json.State.ExitCode,
				Error:      json.State.Error,
				StartedAt:  json.State.StartedAt,
				FinishedAt: json.State.FinishedAt,
				Health:     json.State.Health,
			},
			Started: json.State.Running,
			//todo :may add net config
		}

		//fmt.Println("now in pods " + pod.Status.ContainerStatuses[1].Name)
		//fmt.Println(pod.Status.ContainerStatuses[1].State.Running)
	}
	//pod.Status.ContainerStatuses = containerStatus

}
