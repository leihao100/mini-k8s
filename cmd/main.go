package main

import (
	"MiniK8S/pkg/kubelet"
	"MiniK8S/pkg/util/config/containerConfig"
	"MiniK8S/pkg/util/config/podConfig"
	"MiniK8S/pkg/util/meta"
	"MiniK8S/pkg/util/spec"
	"MiniK8S/pkg/util/status"
)

func main() {
	//cl, err := client.NewClientWithOpts(client.WithVersion("1.43"))
	//if err != nil {
	//	fmt.Println("Unable to create docker client")
	//	panic(err)
	//}
	co := containerConfig.ContainerConfig{
		Name:         "helloworld",
		Cmd:          nil,
		Entrypoint:   nil,
		Env:          nil,
		Image:        "hello-world:latest",
		Volumes:      nil,
		Labels:       nil,
		PortBindings: nil,
		VolumesFrom:  nil,
		Binds:        nil,
		NetworkMode:  "",
		CPULimit:     0,
		MemLimit:     0,
	}
	co1 := containerConfig.ContainerConfig{
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
	var containers []containerConfig.ContainerConfig
	containers = append(containers, co)
	containers = append(containers, co1)
	pod := podConfig.PodConfig{
		ApiVersion: "",
		Kind:       "pod",
		Meta: meta.ObjectMeta{
			Name:      "try",
			Namespace: "try",
		},
		Spec: spec.PodSpec{
			Containers: containers,
		},
		Status: status.PodStatus{},
	}

	k := kubelet.Kubelet{}
	k.Run()
	k.MakePod(pod)

}
