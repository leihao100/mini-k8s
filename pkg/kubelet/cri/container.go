package cri

import (
	"MiniK8S/pkg/api/config"
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
	_ "github.com/docker/docker/pkg/stdcopy"
	"io"
	"os"
)

func GetClient() (Client, error) {
	cil, err := client.NewClientWithOpts(client.WithVersion("1.41"))
	if err != nil {
		panic(err)
		return nil, err
	}
	return &DockerClient{Client: cil}, nil
}

type DockerClient struct {
	Client *client.Client
}

func (c *DockerClient) CreatePause(config config.Container, name string) (*container.CreateResponse, error) {
	ctx := context.Background()
	cl := c.Client
	containerRepoTag := config.Image
	exist := false
	list, err := cl.ImageList(context.Background(), image.ListOptions{})
	for _, repoTag := range list {
		if len(repoTag.RepoTags) == 0 {
			continue
		}
		if repoTag.RepoTags[0] == containerRepoTag {
			exist = true
		}
	}
	if !exist {
		fmt.Println("pulling image ", containerRepoTag)
		res, err := cl.ImagePull(ctx, containerRepoTag, image.PullOptions{})
		//<-ctx.Done()
		if err != nil {
			fmt.Println("Failed to pull image " + containerRepoTag)
			panic(err)
			return nil, err
		}
		res.Close()
	}
	if err != nil {
		fmt.Println("Unable to pull docker client")
		panic(err)
		return nil, err
	}

	var resp container.CreateResponse
	//cl, err = client.NewClientWithOpts(client.WithVersion("1.43"))
	resp, err = cl.ContainerCreate(ctx, &container.Config{
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
		//Mounts:
	}, nil, nil, name)
	if err != nil {
		fmt.Println("Unable to create docker container")
		panic(err)
		return nil, err
	}
	return &resp, nil
}

func (c *DockerClient) CreateContainer(config config.Container, name string) (*container.CreateResponse, error) {
	ctx := context.Background()
	//cl, err := client.NewClientWithOpts(client.WithVersion("1.43"), client.FromEnv, client.WithHost())
	cl := c.Client
	//cl := c
	containerRepoTag := config.Image
	exist := false
	list, err := cl.ImageList(context.Background(), image.ListOptions{})
	for _, repoTag := range list {
		if repoTag.RepoTags[0] == containerRepoTag {
			exist = true
		}
	}
	//config.Binds = append(config.Binds, "/etc:/etc")
	if !exist {
		ct := context.Background()
		fmt.Println("pulling image ", containerRepoTag)
		out, err := cl.ImagePull(ct, containerRepoTag, image.PullOptions{})
		if err != nil {
			fmt.Println("Failed to pull image " + containerRepoTag)
			panic(err)
			return nil, err
		}
		_, err = io.Copy(io.Discard, out)
		out.Close()

	}

	if err != nil {
		fmt.Println("Unable to create docker client")
		panic(err)
		return nil, err
	}
	//pauseId := cl.containerId

	var resp container.CreateResponse
	resp, err = cl.ContainerCreate(ctx, &container.Config{
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
		NetworkMode:  container.NetworkMode("container:" + config.Pause),
		Cgroup:       container.CgroupSpec("container:" + config.Pause),
		Mounts:       c.BuildMount(&config),
	}, nil, nil, name)
	if err != nil {
		fmt.Println("Unable to create docker container")
		panic(err)
		return nil, err
	}
	return &resp, nil

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
		fmt.Println("Unable to stop docker container")
		panic(err)
		return false, err
	}
	return true, nil
}
func (c *DockerClient) ContainerStatus(id string) (types.ContainerJSON, error) {
	ctx := context.Background()
	resp, err := c.Client.ContainerInspect(ctx, id)
	if err != nil {
		return resp, err
	}

	return resp, err
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

func (c *DockerClient) VolumeCreate(v config.Volume) error {
	ctx := context.Background()
	c.Client.VolumeCreate(ctx, volume.CreateOptions{
		Driver: "local",
		Name:   v.Name,
	})
	return nil
}

func (c *DockerClient) BuildMount(con *config.Container) []mount.Mount {
	mnt := make([]mount.Mount, 0)
	for _, m := range con.VolumeMount {
		mnt = append(mnt, mount.Mount{
			Type:   mount.TypeVolume,
			Source: m.Name,
			Target: m.MountPath,
		})
	}
	mnt = append(mnt, mount.Mount{
		Type:   mount.TypeBind,
		Source: "/etc",
		Target: "/etc",
		BindOptions: &mount.BindOptions{
			Propagation: mount.PropagationShared,
		},
	})

	return mnt
}

func (c *DockerClient) Execute(id string, cmd []string) {
	ctx := context.Background()
	idResponse, err := c.Client.ContainerExecCreate(ctx, id, types.ExecConfig{
		Cmd: cmd,
	})
	if err != nil {
		return
	}
	response, err := c.Client.ContainerExecAttach(ctx, idResponse.ID, types.ExecStartCheck{})
	if err != nil {
		return
	}
	defer response.Close()
	io.Copy(os.Stdout, response.Reader)
}
