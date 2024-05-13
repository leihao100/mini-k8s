package main

import (
	"MiniK8S/pkg/api/config"
	"MiniK8S/pkg/api/meta"
	"MiniK8S/pkg/api/status"
	"MiniK8S/pkg/kubelet"
)

func main() {
	//cl, err := client.NewClientWithOpts(client.WithVersion("1.43"))
	//if err != nil {
	//	fmt.Println("Unable to create docker client")
	//	panic(err)
	//}
	co := config.Container{
		Name:         "helloworld",
		Cmd:          nil,
		Entrypoint:   nil,
		Env:          nil,
		Image:        "mysql:latest",
		Volumes:      nil,
		Labels:       nil,
		PortBindings: nil,
		VolumesFrom:  nil,
		Binds:        nil,
		NetworkMode:  "",
		CPULimit:     0,
		MemLimit:     0,
	}
	co1 := config.Container{
		Name:         "nginx",
		Cmd:          nil,
		Entrypoint:   nil,
		Env:          nil,
		Image:        "nginx:latest",
		Volumes:      nil,
		Labels:       nil,
		PortBindings: nil,
		VolumesFrom:  nil,
		Binds:        nil,
		NetworkMode:  "",
		CPULimit:     0,
		MemLimit:     0,
	}
	//var cli cri.Client
	//cli, _ = cri.GetClient()
	//_, err = cli.CreateContainer(co, "hello-world")
	//if err != nil {
	//	panic(err)
	//	fmt.Println("Unable to create docker container")
	//}
	//cli.StartContainer("4340fb1f47f3137793238029acb50b4ab82b0d9c915abd9481c2dd26e18ba8d8")

	//list, err := cl.ImageList(context.Background(), image.ListOptions{})
	//if err != nil {
	//	panic(err)
	//}
	//for _, img := range list {
	//	fmt.Println("a img", img.RepoTags)
	//	if img.RepoTags[0] == "hello-world:latest" {
	//		fmt.Println("hello-world")
	//	}
	//}
	var containers []config.Container
	containers = append(containers, co)
	containers = append(containers, co1)
	pod := config.Pod{
		ApiVersion: "",
		Kind:       "pod",
		Metadata: meta.ObjectMeta{
			Name:      "try",
			Namespace: "try",
		},
		Spec: config.PodSpec{
			Containers: containers,
		},
		Status: status.PodStatus{},
	}

	k := kubelet.Kubelet{}
	k.Run()
	k.MakePod(&pod)
	k.UpdatePodStatusByID(pod.Metadata.Uid)
	//fmt.Println((pod.Status.ContainerStatuses)[1].State.Running)
	//pods := k.GetPods()
	//for _, pod := range pods {
	//	fmt.Println("pod id is", pod.Meta.Uid)
	//	fmt.Println("first container is", pod.Status.ContainerStatuses[0].Name, pod.Status.ContainerStatuses[0].ContainerID)
	//	cl, _ := client.NewClientWithOpts(client.WithVersion("1.43"))
	//	json, err := cl.ContainerInspect(context.Background(), pod.Status.ContainerStatuses[1].ContainerID)
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//	fmt.Println(json)
	//
	//}
}
