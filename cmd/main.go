package main

import (
	"MiniK8S/pkg/kubelet/cri"
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func main() {
	cl, err := client.NewClientWithOpts(client.WithVersion("1.43"))
	if err != nil {
		fmt.Println("Unable to create docker client")
		panic(err)
	}
	//co := containerConfig.ContainerConfig{
	//	Cmd:             nil,
	//	Entrypoint:      nil,
	//	Env:             nil,
	//	Image:           "hello-world",
	//	Volumes:         nil,
	//	Labels:          nil,
	//	ImagePullPolicy: "",
	//	PortBindings:    nil,
	//	VolumesFrom:     nil,
	//	Binds:           nil,
	//	NetworkMode:     "",
	//	CPULimit:        0,
	//	MemLimit:        0,
	//}
	var cli cri.Client
	cli, _ = cri.GetClient()
	//_, err = cli.CreateContainer(co, "hello-world")
	//if err != nil {
	//	panic(err)
	//	fmt.Println("Unable to create docker container")
	//}
	cli.StartContainer("4340fb1f47f3137793238029acb50b4ab82b0d9c915abd9481c2dd26e18ba8d8")

	fmt.Println(cl.ImageList(context.Background(), types.ImageListOptions{}))

}
