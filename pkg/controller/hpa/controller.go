package hpa

import (
	"MiniK8S/pkg/api/config"
	"MiniK8S/pkg/apiClient"
	"MiniK8S/pkg/controller/cache"
)

type hpaController struct {
	podClient    *apiClient.Client
	hpaClient    *apiClient.Client
	deployClient *apiClient.Client
	podInformer  *cache.Informer
	hpaInformer  *cache.Informer
	dpInformer   *cache.Informer
	queue        *cache.WorkQueue
}

func NewController(pi *cache.Informer, hi *cache.Informer, pc *apiClient.Client, hc *apiClient.Client, dc *apiClient.Client) *hpaController {
	hpc := &hpaController{
		hpaClient:    hc,
		podClient:    pc,
		deployClient: dc,
		podInformer:  pi,
		hpaInformer:  hi,
		queue:        cache.NewWorkQueue(),
	}
	hpc.hpaInformer.AddEventHandler(cache.EventHandlerFuncs{
		AddFunc:    hpc.AddHpa,
		UpdateFunc: hpc.UpdateHpa,
		DeleteFunc: hpc.DeleteHpa,
	})
	hpc.podInformer.AddEventHandler(cache.EventHandlerFuncs{
		AddFunc:    hpc.AddPod,
		UpdateFunc: hpc.UpdatePod,
		DeleteFunc: hpc.DeletePod,
	})
	return hpc
}

func (hpc *hpaController) AddPod(obj interface{}) {

}
func (hpc *hpaController) UpdatePod(oldObj, newObj interface{}) {

}
func (hpc *hpaController) DeletePod(obj interface{}) {

}
func (hpc *hpaController) AddHpa(obj interface{}) {

}
func (hpc *hpaController) UpdateHpa(oldObj, newObj interface{}) {

}
func (hpc *hpaController) DeleteHpa(obj interface{}) {

}

func (hpc *hpaController) CalculateTarget(hpa *config.HorizontalPodAutoscaler) {
	//resource :=
}
