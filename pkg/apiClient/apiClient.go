package apiClient

import (
	"MiniK8S/config"
	core "MiniK8S/pkg/api/config"
	"MiniK8S/pkg/api/types"
	"MiniK8S/pkg/api/url"
	"MiniK8S/pkg/api/watch"
	"bytes"
	"errors"
	"io"
	"log"
	"net/http"
	"time"
)

type RequestType string

const (
	Create RequestType = "create"
	Delete RequestType = "delete"
	Get    RequestType = "get"
	Watch  RequestType = "watch"
	Status RequestType = "status"
)

const HttpStatusNotSend = 0
const ReconnectInterval = 5

// Client is a REST client for Kubernetes API
type Client struct {
	// ApiServerUrl is the URL of the API server
	ApiServerUrl string
	// MiddleURL is a URL generator
	MiddleURL    url.URL
	ResourceType types.ApiObjectType
}

func NewRESTClient(ty types.ApiObjectType) *Client {
	newURL := url.URL{}
	field := ty
	switch ty {
	case types.PodObjectType:
		field = "pods"
	case types.ServiceObjectType:
		field = "services"
	case types.DeploymentObjectType:
		field = "deployments"
	case types.HeartbeatObjectType:
		field = "heartbeats"
	case types.HorizontalPodAutoscalerObjectType:
		field = "hpas"
	case types.NodeObjectType:
		field = "nodes"
	case types.DnsObjectType:
		field = "dnss"
	case types.PersistentVolumeObjectType:
		field = "persistentVolumes"
	case types.PersistentVolumeClaimObjectType:
		field = "persistentVolumeClaims"
	case types.StorageClassObjectType:
		field = "storageClasses"
	}
	newURL.Init("v1", string(field))
	apiserverURL := config.ApiServerHost() + config.ApiServerPort()
	return &Client{
		ApiServerUrl: apiserverURL,
		MiddleURL:    newURL,
		ResourceType: ty,
	}
}

func putBytes(URL string, content []byte) (*http.Response, error) {
	cli := &http.Client{}
	req, err := http.NewRequest(http.MethodPut, URL, bytes.NewReader(content))
	if err != nil {
		log.Println("[utils][http][PutBytes] http.NewRequest create failed", err)
		return nil, err
	}
	return cli.Do(req)
}

func (c *Client) PutObject(object core.ApiObject) (int, error) {
	putUrl := c.BuildFullURL(Create, "")
	content, err := object.JsonMarshal()
	if err != nil {
		log.Println("[Client] http.Put JsonMarshal failed", err)
		return HttpStatusNotSend, err
	}

	resp, err := putBytes(putUrl, content)
	if err != nil {
		log.Println("[Client] http.Put failed", err)
		return HttpStatusNotSend, err
	}

	if resp.StatusCode == http.StatusOK {
		return resp.StatusCode, nil
	} else {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println("[Client] http.Put StatusCode not http.StatusOK, http.Put io.ReadAll failed", err)
		} else {
			log.Println("[Client] http.Put StatusCode not http.StatusOK, ", string(body))
		}
		return resp.StatusCode, errors.New("StatusCode not 200")
	}
}

func (c *Client) URL() string {
	return c.ApiServerUrl
}

func (c *Client) BuildURL(requestType RequestType) string {
	res := c.ApiServerUrl
	switch requestType {
	case Create:
		return res + c.MiddleURL.CreateURL()
	case Delete:
		return res + c.MiddleURL.DeleteURL()
	case Get:
		return res + c.MiddleURL.GetURL()
	case Watch:
		return res + c.MiddleURL.WatchURL()
	case Status:
		return res + c.MiddleURL.StatusURL()
	default:
		return ""
	}
}

// BuildFullURL builds a full URL for a given request type and resource name
func (c *Client) BuildFullURL(requestType RequestType, resourceName string) string {
	res := c.ApiServerUrl
	switch requestType {
	case Create:
		res = res + c.MiddleURL.CreateURL()
	case Delete:
		res = res + c.MiddleURL.DeleteURL()
	case Get:
		res = res + c.MiddleURL.GetURL()
	case Watch:
		res = res + c.MiddleURL.WatchURL()
	case Status:
		res = res + c.MiddleURL.StatusURL()
	default:
		return ""
	}
	if resourceName != "" {
		res = res + "/" + resourceName
	}
	return res
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

func (c *Client) GetObject(name string) (int, core.ApiObject, error) {
	getUrl := c.BuildFullURL(Get, name)
	resp, err := http.Get(getUrl)
	if err != nil {
		log.Println("[Client] http.Get failed", err)
		return HttpStatusNotSend, nil, err
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("[Client] http.Get response io.ReadAll failed", err)
		return resp.StatusCode, nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return resp.StatusCode, nil, errors.New("StatusCode not 200, response code: " + err.Error() + "response body: " + string(content))
	}

	object := core.NewApiObject(c.ResourceType)
	err = object.JsonUnmarshal(content)
	if err != nil {
		log.Println("[Client] http.Get response json.Unmarshal failed", err)
		return resp.StatusCode, nil, err
	}

	return resp.StatusCode, object, nil
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

func (c *Client) GetAll() (objectList core.ApiObjectList, err error) {
	getUrl := c.BuildFullURL(Get, "")
	resp, err := http.Get(getUrl)
	if err != nil {
		log.Println("[RESTClient] WatchAll Failed: ", err)
		return nil, err
	}
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("[RESTClient] http.GetAll response io.ReadAll(resp.Body) failed", err)
		return nil, err
	}
	objectList = core.NewApiObjectList(c.ResourceType)
	if len(content) == 0 {
		return objectList, nil
	}
	err = objectList.JsonUnmarshal(content)
	if err != nil {
		log.Printf("[RESTClient] http.Get response json.Unmarshal failed, err %v\n", err)
		return nil, err
	}
	return objectList, nil
}

func (c *Client) WatchAll() (watch.Interface, error) {
	watchUrl := c.BuildFullURL(Watch, "")
	resp, err := http.Get(watchUrl)

	if err != nil {
		log.Println("[RESTClient] WatchAll Failed: ", err)
		// sleep some time before retry
		time.Sleep(time.Second * time.Duration(ReconnectInterval))
		return nil, err
	}

	decoder := watch.NewEtcdEventDecoder(resp.Body, c.ResourceType)
	streamWatcher := watch.NewStreamWatcher(decoder)

	return streamWatcher, nil
}

func (c *Client) Watch(name string) (watch.Interface, error) {
	watchURL := c.BuildFullURL(Watch, name)
	resp, err := http.Get(watchURL)

	if err != nil {
		log.Printf("[RESTClient] Watch %v %v Failed: %v\n", c.ResourceType, name, err)
		// sleep some time before retry
		time.Sleep(time.Second * time.Duration(ReconnectInterval))
		return nil, err
	}

	log.Printf("[RESTClient] Watch %v %v start\n", c.ResourceType, name)

	decoder := watch.NewEtcdEventDecoder(resp.Body, c.ResourceType)
	streamWatcher := watch.NewStreamWatcher(decoder)

	return streamWatcher, nil
}
