package deployment

import (
	"MiniK8S/pkg/api/config"
	"MiniK8S/pkg/api/selector"
	"MiniK8S/pkg/apiClient"
	"MiniK8S/pkg/controller/cache"
	"context"
	"github.com/google/uuid"
	"reflect"
	"time"
)

type DeploymentController struct {
	podClient       *apiClient.Client
	deployClient    *apiClient.Client
	podInformer     *cache.Informer
	replicaInformer *cache.Informer
	replicaMap      map[uuid.UUID]*config.Deployment
	pods            map[uuid.UUID][]*config.Pod
	queue           *cache.WorkQueue
}

func NewController(pi *cache.Informer, ri *cache.Informer, pc *apiClient.Client, dc *apiClient.Client) *DeploymentController {
	dpc := &DeploymentController{
		deployClient:    dc,
		podClient:       pc,
		podInformer:     pi,
		replicaInformer: ri,
		queue:           cache.NewWorkQueue(),

		replicaMap: make(map[uuid.UUID]*config.Deployment),
	}
	dpc.replicaInformer.AddEventHandler(cache.EventHandlerFuncs{
		AddFunc:    dpc.AddDeployment,
		UpdateFunc: dpc.UpdateDeployment,
		DeleteFunc: dpc.DeleteDeployment,
	})
	dpc.podInformer.AddEventHandler(cache.EventHandlerFuncs{
		AddFunc:    dpc.AddPod,
		UpdateFunc: dpc.UpdatePod,
		DeleteFunc: dpc.DeletePod,
	})
	return dpc
}

func (dpc *DeploymentController) AddDeployment(obj interface{}) {
	dp := obj.(config.Deployment)
	dpc.queue.Add(dp)
}
func (dpc *DeploymentController) DeleteDeployment(obj interface{}) {

}
func (dpc *DeploymentController) UpdateDeployment(oldObj, newObj interface{}) {
	dp := newObj.(config.Deployment)
	dpc.queue.Add(dp)
}
func (dpc *DeploymentController) AddPod(obj interface{}) {
	pd := obj.(config.Pod)
	dps := dpc.GetDpsByPod(&pd)
	for _, dp := range dps {
		dpc.queue.Add(dp)
	}
	//dpc.queue.Add(pd)
}
func (dpc *DeploymentController) DeletePod(obj interface{}) {
	//pd := obj.(config.Pod)
	//owners := pd.Metadata.OwnerReferences
	//for _, owner := range owners {
	//	o := owner.UID
	//	dpc.queue.Add(dpc.replicaMap[o])
	//}
	pd := obj.(config.Pod)
	dps := dpc.GetDpsByPod(&pd)
	for _, dp := range dps {
		dpc.queue.Add(dp)
	}
}

func (dpc *DeploymentController) UpdatePod(oldObj, newObj interface{}) {
	//pd := newObj.(config.Pod)
	//dpc.queue.Add(pd)
	oldpd := oldObj.(config.Pod)
	newpd := newObj.(config.Pod)
	if reflect.DeepEqual(oldpd.Metadata.Labels, newpd.Metadata.Labels) {
		return
	}

	olddps := dpc.GetDpsByPod(&oldpd)
	newdps := dpc.GetDpsByPod(&newpd)

	for _, dp := range olddps {
		dpc.queue.Add(dp)
	}
	for _, dp := range newdps {
		dpc.queue.Add(dp)
	}
}
func (dc *DeploymentController) Run(ctx context.Context, cancel context.CancelFunc) {
	go func() {
		defer cancel()
		for true {
			select {
			case <-ctx.Done():
				return
			default:
				obj, ok := dc.queue.Get()
				if !ok { //此时队列为空
					time.Sleep(1 * time.Second)
				}
				dp := obj.(config.Deployment)
				dc.Sync(&dp)
			}
		}
	}()
}

func (dc *DeploymentController) ListDeployments() []*config.Deployment {
	result := make([]*config.Deployment, 0)
	for _, dp := range dc.replicaMap {
		result = append(result, dp)
	}
	return result
}

func (dc *DeploymentController) Sync(dp *config.Deployment) {

}

func (dc *DeploymentController) SelectDpByLabelSelector(labelSelector selector.LabelSelector) []*config.Deployment {
	dps := dc.replicaInformer.List()
	var result []*config.Deployment
	for _, dp := range dps {
		actualDp := dp.(*config.Deployment)
		if selector.LabelSelectorCompare(labelSelector, actualDp.Metadata.Labels) {
			result = append(result, actualDp)
		}
	}
	return result
}

func (dc *DeploymentController) GetDpsByPod(pod *config.Pod) []*config.Deployment {
	dps := dc.replicaInformer.List()
	var result []*config.Deployment
	for _, dp := range dps {
		actualDp := dp.(*config.Deployment)
		if selector.LabelCompare(actualDp.Metadata.Labels, pod.Metadata.Labels) {
			result = append(result, actualDp)
		}
	}
	return result
}
