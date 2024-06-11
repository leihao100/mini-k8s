package main

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
)

func main() {
	cli, err := client.NewClientWithOpts(client.WithVersion("1.41"))
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	// 创建两个 Volume
	volume1, err := cli.VolumeCreate(ctx, volume.CreateOptions{Driver: "local", Name: "volume1"})
	if err != nil {
		panic(err)
	}

	volume2, err := cli.VolumeCreate(ctx, volume.CreateOptions{Driver: "local", Name: "volume2"})
	if err != nil {
		panic(err)
	}

	// 创建第一个容器并挂载 Volume
	container1, err := cli.ContainerCreate(
		ctx,
		&container.Config{
			Image: "ubuntu",
			Cmd:   []string{"tail", "-f", "/dev/null"},
		},
		&container.HostConfig{
			Mounts: []mount.Mount{
				{
					Type:   mount.TypeVolume,
					Source: volume1.Name,
					Target: "/tmp",
				},
				{
					Type:   mount.TypeVolume,
					Source: volume2.Name,
					Target: "/etc",
				},
			},
		},
		nil,
		nil,
		"container1",
	)
	if err != nil {
		panic(err)
	}

	// 创建第二个容器并挂载 Volume
	container2, err := cli.ContainerCreate(
		ctx,
		&container.Config{
			Image: "ubuntu",
			Cmd:   []string{"tail", "-f", "/dev/null"},
		},
		&container.HostConfig{
			Mounts: []mount.Mount{
				{
					Type:   mount.TypeVolume,
					Source: volume1.Name,
					Target: "/root/tmp",
				},
				{
					Type:   mount.TypeVolume,
					Source: volume2.Name,
					Target: "/root/etc",
				},
			},
		},
		nil,
		nil,
		"container2",
	)
	if err != nil {
		panic(err)
	}

	// 启动两个容器
	if err := cli.ContainerStart(ctx, container1.ID, container.StartOptions{}); err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, container2.ID, container.StartOptions{}); err != nil {
		panic(err)
	}

	fmt.Println("Container1 ID:", container1.ID)
	fmt.Println("Container2 ID:", container2.ID)
}
