package container

import (
	"MiniK8S/pkg/util/config"
	"context"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	_ "github.com/docker/docker/pkg/stdcopy"
)

func createContainer(config config.ContainerConfig, name string) (bool, error) {
	ctx := context.Background()
	//cl, err := client.NewClientWithOpts(client.WithVersion("1.43"), client.FromEnv, client.WithHost())
	cl, err := client.NewClientWithOpts(client.WithVersion("1.43"))
	if err != nil {
		fmt.Println("Unable to create docker client")
		panic(err)
		return false, err
	}

	_, err = cl.ContainerCreate(ctx, &container.Config{
		Image:      config.Image,
		Cmd:        config.Cmd,
		Env:        config.Env,
		Entrypoint: config.Entrypoint,
		Volumes:    config.Volumes,
		Labels:     config.Labels,
	}, &container.HostConfig{
		Binds:        config.Binds,
		PortBindings: config.PortBindings,
		VolumesFrom:  config.VolumesFrom,
		NetworkMode:  container.NetworkMode(config.NetworkMode),
	}, nil, nil, name)
	if err != nil {
		fmt.Println("Unable to create docker container")
		panic(err)
		return false, err
	}
	return true, nil

}
func startContainer(name string) (bool, error) {
	ctx := context.Background()
	//cl, err := client.NewClientWithOpts(client.WithVersion("1.43"), client.FromEnv, client.WithHost())
	cl, err := client.NewClientWithOpts(client.WithVersion("1.43"))
	if err != nil {
		fmt.Println("Unable to create docker client")
		panic(err)
		return false, err
	}
	err = cl.ContainerStart(ctx, name, container.StartOptions{})
	if err != nil {
		fmt.Println("Unable to start docker container")
		panic(err)
		return false, err
	}
	return true, nil
}

func stopContainer(name string) (bool, error) {
	ctx := context.Background()
	//cl, err := client.NewClientWithOpts(client.WithVersion("1.43"), client.FromEnv, client.WithHost())
	cl, err := client.NewClientWithOpts(client.WithVersion("1.43"))
	if err != nil {
		fmt.Println("Unable to create docker client")
		panic(err)
		return false, err
	}
	err = cl.ContainerStop(ctx, name, container.StopOptions{})
	if err != nil {
		fmt.Println("Unable to start docker container")
		panic(err)
		return false, err
	}
	return true, nil
}

//func startPod() {
//	ctx := context.Background()
//	cl, err := client.NewClientWithOpts(client.WithVersion("1.43"))
//	if err != nil {
//		fmt.Println("Unable to create docker client")
//		panic(err)
//	}
//
//}
