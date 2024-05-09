package cri

import (
	"MiniK8S/pkg/util/config"
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	_ "github.com/docker/docker/pkg/stdcopy"
)

func GetClient() (Client, error) {
	cil, err := client.NewClientWithOpts(client.WithVersion("1.43"))
	if err != nil {
		panic(err)
		return nil, err
	}
	return &DockerClient{Client: cil}, nil
}

type DockerClient struct {
	Client *client.Client
}

func (c *DockerClient) CreateContainer(config config.ContainerConfig, name string) (bool, error) {
	ctx := context.Background()
	//cl, err := client.NewClientWithOpts(client.WithVersion("1.43"), client.FromEnv, client.WithHost())
	cl, err := client.NewClientWithOpts(client.WithVersion("1.43"))
	//cl := c
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
func (c *DockerClient) StartContainer(id string) (bool, error) {
	ctx := context.Background()
	err := c.Client.ContainerStart(ctx, id, container.StartOptions{})
	if err != nil {
		fmt.Println("Unable to start docker container")
		panic(err)
		return false, err
	}
	return true, nil
}

func (c *DockerClient) StopContainer(id string) (bool, error) {
	ctx := context.Background()

	err := c.Client.ContainerStop(ctx, id, container.StopOptions{})
	if err != nil {
		fmt.Println("Unable to start docker container")
		panic(err)
		return false, err
	}
	return true, nil
}
func (c *DockerClient) ContainerStatus(id string) (bool, int, error) {
	ctx := context.Background()
	resp, err := c.Client.ContainerInspect(ctx, id)
	if err != nil {
		return false, 0, err
	}
	return resp.State.Running, resp.State.ExitCode, nil
}
func (c *DockerClient) RemoveContainer(id string) error {
	ctx := context.Background()
	err := c.Client.ContainerRemove(ctx, id, container.RemoveOptions{})
	if err != nil {
		fmt.Println("Unable to remove docker container")
		panic(err)
		return err
	}
	return nil
}

func (c *DockerClient) Close() error {
	return c.Client.Close()
}

func (c *DockerClient) ListContainers() []types.Container {
	containers, err := c.Client.ContainerList(context.Background(), container.ListOptions{})
	if err != nil {
		panic(err)
		return nil
	}
	return containers
}
