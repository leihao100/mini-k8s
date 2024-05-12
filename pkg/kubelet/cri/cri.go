package cri

import (
	"MiniK8S/pkg/api/config"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
)

type Client interface {
	CreateContainer(Config config.Container, name string) (*container.CreateResponse, error)
	StartContainer(id string) (bool, error)
	StopContainer(id string) (bool, error)
	ContainerStatus(id string) (types.ContainerJSON, error)
	RemoveContainer(id string) error
	Close() error
	ListContainers() []types.Container
	CreatePause(Config config.Container, name string) (*container.CreateResponse, error)
}
