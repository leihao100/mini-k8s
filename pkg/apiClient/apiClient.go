package apiClient

import (
	"MiniK8S/config"
	"MiniK8S/pkg/api/types"
	"MiniK8S/pkg/api/url"
	"bytes"
	"io"
	"net/http"
)

type RequestType string

const (
	Create RequestType = "create"
	Delete RequestType = "delete"
	Get    RequestType = "get"
	Watch  RequestType = "watch"
	Status RequestType = "status"
)

type Client struct {
	ApiServerUrl string
	ResourceUrl  url.URL
	ResourceType types.ApiObjectType
}

func NewRESTClient(ty types.ApiObjectType) *Client {
	newURL := url.URL{}
	field := ty
	newURL.Init("v1", string(field))
	apiserverURL := config.ApiServerHost() + config.ApiServerPort()
	return &Client{
		ApiServerUrl: apiserverURL,
		ResourceUrl:  newURL,
		ResourceType: ty,
	}
}

func (c *Client) URL() string {
	return c.ApiServerUrl
}

func (c *Client) BuildURL(requestType RequestType) string {
	res := c.ApiServerUrl
	switch requestType {
	case Create:
		return res + c.ResourceUrl.CreateURL()
	case Delete:
		return res + c.ResourceUrl.DeleteURL()
	case Get:
		return res + c.ResourceUrl.GetURL()
	case Watch:
		return res + c.ResourceUrl.WatchURL()
	case Status:
		return res + c.ResourceUrl.StatusURL()
	default:
		return ""
	}
}

func (c *Client) Post(resourceURL string, context []byte) io.ReadCloser {
	postUrl := resourceURL
	cli := &http.Client{}
	//创建请求
	req, err := http.NewRequest(http.MethodPost, postUrl, bytes.NewReader(context))
	if err != nil {
		panic(err)
	}
	//发送请求
	response, err := cli.Do(req)
	if err != nil {
		panic(err)
	}
	//检查返回
	/*//defer func(Body io.ReadCloser) {
	//	err := Body.Close()
	//	if err != nil {
	//		panic(err)
	//	}
	//}(response.Body)
	*/
	//response即为返回
	return response.Body
}

func (c *Client) Get(resourceURL string, context []byte) io.ReadCloser {
	getUrl := resourceURL
	cli := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, getUrl, bytes.NewReader(context))
	if err != nil {
		panic(err)
	}
	response, err := cli.Do(req)
	if err != nil {
		panic(err)
	}

	return response.Body
}

func (c *Client) Put(resourceURL string, context []byte) io.ReadCloser {
	putUrl := resourceURL
	cli := &http.Client{}
	req, err := http.NewRequest(http.MethodPut, putUrl, bytes.NewReader(context))
	if err != nil {
		panic(err)
	}
	response, err := cli.Do(req)
	if err != nil {
		panic(err)
	}

	return response.Body
}

func (c *Client) Delete(resourceURL string, context []byte) io.ReadCloser {
	deleteUrl := resourceURL
	cli := &http.Client{}
	req, err := http.NewRequest(http.MethodDelete, deleteUrl, bytes.NewReader(context))
	if err != nil {
		panic(err)
	}
	response, err := cli.Do(req)
	if err != nil {
		panic(err)
	}

	return response.Body
}
