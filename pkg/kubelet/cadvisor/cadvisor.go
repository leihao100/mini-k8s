package cadvisor

import "github.com/google/cadvisor/client"

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
