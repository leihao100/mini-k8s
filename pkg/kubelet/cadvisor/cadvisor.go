package cadvisor

import (
	"github.com/google/cadvisor/client"
	v1 "github.com/google/cadvisor/info/v1"
)

const CAdvisorURL = "http://localhost:9000"

type CAdvisorClient struct {
	cli client.Client
}

func NewCAdvisor(url string) *CAdvisorClient {
	cli, err := client.NewClient(url)
	if err != nil {
		panic(err)
	}
	return &CAdvisorClient{
		cli: *cli,
	}
}

func (c *CAdvisorClient) Start() error {
	return nil
}

func (c *CAdvisorClient) Inspect(query *v1.ContainerInfoRequest) ([]v1.ContainerInfo, error) {
	return c.cli.AllDockerContainers(query)

}
