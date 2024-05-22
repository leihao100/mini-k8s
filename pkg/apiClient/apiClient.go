package apiClient

import "MiniK8S/pkg/api/types"

type Client struct {
	ApiServerUrl string
	ResourceUrl  string
	ResourceType types.ApiObjectType
}

func NewRESTClient() *Client {
	return &Client{
		ApiServerUrl: "",
		ResourceUrl:  "",
	}
}

func Post(client Client, resourceType types.ApiObjectType) {

}
