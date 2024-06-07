package prometheus

import (
	"MiniK8S/pkg/api/config"
	"MiniK8S/pkg/apiClient"
	"MiniK8S/pkg/controller/cache"
	"MiniK8S/utils/net"
	"fmt"
	"strings"
)

type PrometheusController struct {
	podClient    *apiClient.Client
	nodeClient   *apiClient.Client
	podInformer  *cache.Informer
	nodeInformer *cache.Informer
	queue        *cache.WorkQueue
}

func NewPrometheusController(pc, nc *apiClient.Client, pi, ni *cache.Informer) *PrometheusController {
	ptc := &PrometheusController{
		podClient:    pc,
		nodeClient:   nc,
		podInformer:  pi,
		nodeInformer: ni,
	}
	ptc.nodeInformer.AddEventHandler(cache.EventHandlerFuncs{
		AddFunc:    ptc.AddNode,
		UpdateFunc: ptc.UpdateNode,
		DeleteFunc: ptc.DeleteNode,
	})
	ptc.podInformer.AddEventHandler(cache.EventHandlerFuncs{
		AddFunc:    ptc.AddPod,
		UpdateFunc: ptc.UpdatePod,
		DeleteFunc: ptc.DeletePod,
	})
	return ptc
}
func (ptc *PrometheusController) AddPod(obj interface{}) {
	//pd := obj.(*config.Pod)
}

func (ptc *PrometheusController) UpdatePod(oldobj, newobj interface{}) {

	pd := newobj.(*config.Pod)
	if strings.EqualFold(pd.Metadata.Namespace, "prometheus") {
		if pd.Status.PodIP == "" || pd.Metadata.Generation == 1 {
			return
		} else {
			net.AddPrometheus(pd.Metadata.Name, pd.Status.PodIP+":9090")
			pd.Metadata.Generation = 1
			url := ptc.podClient.BuildURL(apiClient.Create)
			buf, _ := pd.JsonMarshal()
			ptc.podClient.Put(url, buf)
		}
	}
}
func (ptc *PrometheusController) DeletePod(obj interface{}) {
	pd := obj.(*config.Pod)
	if strings.EqualFold(pd.Metadata.Namespace, "prometheus") {
		if pd.Status.PodIP == "" {
			return
		} else {
			net.RemovePrometheus(pd.Metadata.Name)
			//net.AddPrometheus(pd.Metadata.Name, pd.Status.PodIP+":9090")
		}
	}

}

func (ptc *PrometheusController) AddNode(obj interface{}) {
	fmt.Println("[prometheus] Add Node")
	node := obj.(*config.Node)
	net.AddPrometheus(node.Metadata.Name, node.Status.Addresses.Address+":9090")
}
func (ptc *PrometheusController) UpdateNode(obj, oldobj interface{}) {

}
func (ptc *PrometheusController) DeleteNode(obj interface{}) {
	node := obj.(*config.Node)
	net.RemovePrometheus(node.Metadata.Name)
}
