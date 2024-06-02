package kubelet

import (
	"MiniK8S/pkg/api/config"
	"MiniK8S/pkg/api/status"
	apitypes "MiniK8S/pkg/api/types"
	"MiniK8S/pkg/api/watch"
	"MiniK8S/pkg/apiClient"
	"MiniK8S/pkg/apiClient/listwatch"
	"MiniK8S/pkg/kubelet/cri"
	"MiniK8S/pkg/kubelet/pod"
	"context"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/google/uuid"
	"io"
	"reflect"
	"sync"
	"time"
)

const pauseName = "mirrorgooglecontainers/pause:latest"

type Kubelet struct {
	node           config.Node
	cli            cri.Client
	podManager     *pod.PodManager
	podClient      *apiClient.Client
	podListWatcher listwatch.ListerWatcher
	lock           sync.Locker
}

func NewKubelet(node config.Node) *Kubelet {
	cli, _ := cri.GetClient()
	return &Kubelet{
		node:           node,
		cli:            cli,
		podManager:     pod.NewPodManager(),
		podClient:      apiClient.NewRESTClient(apitypes.PodObjectType),
		podListWatcher: nil,
	}
}

func (k *Kubelet) Run(ctx context.Context, cancel context.CancelFunc) error {

	k.podListWatcher = listwatch.NewListWatchFromClient(k.podClient)
	go func() {
		defer cancel()
		k.ListAndWatch(ctx)
	}()
	return nil
}

// just for test
func (k *Kubelet) SendMessage() {

	url := k.podClient.BuildURL(apiClient.Create)
	fmt.Println("my yrl is" + url)
	res := k.podClient.Get(url, nil)
	for {
		body, err := io.ReadAll(res)
		if err != nil {
			panic(err)
		}
		if len(body) == 0 {
			continue
		}
		fmt.Println(string(body))
	}

}

func (k *Kubelet) Stop() {

}

func (k *Kubelet) CreatePodPause(pod *config.Pod) string {
	fmt.Println("[kubelet] CreatePodPause")
	uid := pod.Metadata.Uid.String()

	name := pod.Metadata.Namespace + "_" + pod.Metadata.Name + "_pause_" + uid
	if pod.Metadata.Namespace == "" {
		name = "defaultNameSpace" + name
	}
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
	fmt.Println("[kubelet] makePod" + pod.GetUID().String())
	//pod.Metadata.Uid, _ = uuid.NewUUID()
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
		if pod.Metadata.Namespace == "" {
			containerName = "defaultNameSpace" + containerName
		}
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
	newContainerStatus := status.ContainerStatus{
		Name:         pod.Metadata.Namespace + "_" + pod.Metadata.Name + "_pause_",
		ContainerID:  pauseID,
		ImageID:      "",
		Image:        pauseName,
		State:        types.ContainerState{},
		Started:      true,
		RestartCount: 0,
	}
	pod.Status.ContainerStatuses = append(pod.Status.ContainerStatuses, newContainerStatus)
	k.UpdatePodStatusByID(pod.Metadata.Uid)
	msg, _ := pod.JsonMarshal()
	url := k.podClient.BuildURL(apiClient.Create)
	k.podClient.Put(url, msg)

}

func (k *Kubelet) ModifyPod(pod *config.Pod) {
	fmt.Println("[kubelet] modifyPod" + pod.GetUID().String())
	old := k.podManager.GetPodById(pod.GetUID())
	if old == nil {
		//it is a new pod
		k.MakePod(pod)
		return
	}
	//compare spec
	if reflect.DeepEqual(old.Spec, pod.Spec) {
		return
	}
	k.RemovePod(old)
	k.MakePod(pod)
}

func (k *Kubelet) RemovePod(pod *config.Pod) {
	fmt.Println("[kubelet] RemovePod" + pod.GetUID().String())
	uid := pod.Metadata.Uid
	k.podManager.GetPodById(uid)
	for _, container := range pod.Status.ContainerStatuses {
		_, err := k.cli.StopContainer(container.ContainerID)
		if err != nil {
			panic(err)
		}
		err = k.cli.RemoveContainer(container.ContainerID)
		if err != nil {
			panic(err)
		}
	}
	k.podManager.DeletePodById(uid)

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
	fmt.Println("[kubelet] UpdatePodStatusByID" + pod.Metadata.Name)
	//fmt.Println(len(pod.Status.ContainerStatuses))
	containerStatus := pod.Status.ContainerStatuses
	//pod.Status.PodIP = containerStatus[0].State.

	for i, Status := range containerStatus {
		json, err := k.cli.ContainerStatus(Status.ContainerID)
		if err != nil {
			return
		}
		if i == 0 {
			pod.Status.PodIP = json.NetworkSettings.IPAddress
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

func (k *Kubelet) ListAndWatch(ctx context.Context) {
	podList, err := k.podListWatcher.List(config.ListOptions{
		Kind:            string(apitypes.PodObjectType),
		APIVersion:      "",
		LabelSelector:   "",
		FieldSelector:   "",
		Watch:           false,
		ResourceVersion: "",
		TimeoutSeconds:  nil,
	})
	if err != nil {
		panic(err)
	}
	list := podList.GetItems()
	for _, p := range list {
		fmt.Println("[kubelet] ListAndWatch" + p.GetUID().String())
		pod := p.(*config.Pod)
		k.podManager.AddPod(pod.Metadata.Uid, k.podManager.MakePodName(pod), pod)
	}
	go func(k *Kubelet) {
		for {
			for _, pod := range k.podManager.GetPods() {
				err := k.inspectPod(ctx, pod)
				if err != nil {
					panic(err)
					return
				}
			}
			time.Sleep(5 * time.Second)
		}
	}(k)

	w, err := k.podListWatcher.Watch(config.ListOptions{
		Kind:            string(apitypes.PodObjectType),
		APIVersion:      "",
		LabelSelector:   "",
		FieldSelector:   "",
		Watch:           true,
		ResourceVersion: "",
		TimeoutSeconds:  nil,
	})
	if err != nil {
		panic("kubelet watch pod failed " + err.Error())
	}

	err = k.HandleWatch(w, ctx)
	w.Stop() // stop watch

}

func (k *Kubelet) HandleWatch(w watch.Interface, ctx context.Context) error {

	for {
		select {
		case <-ctx.Done():
			return errors.New("kubelet watch context canceled")
		case event := <-w.ResultChan():
			pod := event.Object.(*config.Pod)
			if pod.Spec.NodeName == k.node.Metadata.Name {
				fmt.Println("[kubelet] Handle Pod Watch Event")
				switch event.Type {
				case watch.Added:
					k.MakePod(pod)
				case watch.Modified:
					k.ModifyPod(pod)
				case watch.Deleted:
					k.RemovePod(pod)
				case watch.Error:
					panic("watch occur error")
				case watch.Bookmark:
				default:
					panic("it should never happen")
				}
			}
		}
	}
}

func (k *Kubelet) inspectPod(ctx context.Context, pod *config.Pod) error {
	fmt.Println("[kubelet] inspectPod")
	old := make([]status.ContainerStatus, 0)
	for _, podstatus := range pod.Status.ContainerStatuses {
		old = append(old, podstatus)
	}
	k.UpdatePodStatusByID(pod.Metadata.Uid)
	phase := status.PodRunning
	for _, podstatus := range pod.Status.ContainerStatuses {
		if podstatus.State.Running == false {
			if podstatus.State.ExitCode != 0 {
				fmt.Println("[kubelet] Pod Status Error: " + podstatus.Name + " Pod failed: " + podstatus.State.Status)
				phase = status.PodFailed
			} else {
				fmt.Println("[kubelet] Pod Status: " + podstatus.Name + " Pod Completed")
				phase = status.PodSucceeded
			}
		}
	}
	if phase != status.PodRunning {
		fmt.Println("[kubelet] pod name is", pod.GetName(), " phase is ", phase)
	}
	if phase != status.PodRunning || !reflect.DeepEqual(old, pod.Status.ContainerStatuses) || pod.Status.PodIP == "" {
		fmt.Println("[kubelet] pod name is", pod.GetName(), " and pod status has changed")
		msg, _ := pod.JsonMarshal()
		url := k.podClient.BuildURL(apiClient.Create)
		k.podClient.Put(url, msg)
		//pod.Status.ContainerStatuses =
	}
	return nil
}
