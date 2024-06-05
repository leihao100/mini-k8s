package deployment

import (
	"MiniK8S/pkg/api/config"
	"MiniK8S/pkg/api/meta"
	"MiniK8S/pkg/api/selector"
	"MiniK8S/pkg/api/status"
	"MiniK8S/pkg/api/types"
	"MiniK8S/pkg/apiClient"
	"MiniK8S/pkg/controller/cache"
	"context"
	"fmt"
	"github.com/docker/docker/testutil"
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
	lastSync        time.Time
}

func NewController(pi *cache.Informer, ri *cache.Informer, pc *apiClient.Client, dc *apiClient.Client) *DeploymentController {
	dpc := &DeploymentController{
		deployClient:    dc,
		podClient:       pc,
		podInformer:     pi,
		replicaInformer: ri,
		queue:           cache.NewWorkQueue(),
		lastSync:        time.Now(),

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
	dp := obj.(*config.Deployment)
	dpc.queue.Add(dp)
}
func (dpc *DeploymentController) DeleteDeployment(obj interface{}) {

}
func (dpc *DeploymentController) UpdateDeployment(oldObj, newObj interface{}) {
	dp := newObj.(*config.Deployment)
	dpc.queue.Add(dp)
}
func (dpc *DeploymentController) AddPod(obj interface{}) {
	pd := obj.(*config.Pod)
	dps := dpc.GetDpsByPod(pd)
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
	fmt.Println("[dpController] delete pod")
	pd := obj.(*config.Pod)
	dps := dpc.GetDpsByPod(pd)
	for _, dp := range dps {
		fmt.Println("[dpController] handle delete pod: Adding dp into work queue : ", dp.GetName())
		dpc.queue.Add(dp)
	}
}

func (dpc *DeploymentController) UpdatePod(oldObj, newObj interface{}) {
	//pd := newObj.(config.Pod)
	//dpc.queue.Add(pd)
	if newObj == nil {
		return
	}
	if oldObj == nil {
		newpd := newObj.(*config.Pod)
		newdps := dpc.GetDpsByPod(newpd)
		for _, dp := range newdps {
			dpc.queue.Add(dp)
		}
	} else {
		oldpd := oldObj.(*config.Pod)
		newpd := newObj.(*config.Pod)
		if reflect.DeepEqual(oldpd.Metadata.Labels, newpd.Metadata.Labels) {
			return
		}

		olddps := dpc.GetDpsByPod(oldpd)
		newdps := dpc.GetDpsByPod(newpd)

		for _, dp := range olddps {
			dpc.queue.Add(dp)
		}
		for _, dp := range newdps {
			dpc.queue.Add(dp)
		}
	}

}
func (dc *DeploymentController) Run(ctx context.Context, cancel context.CancelFunc) {
	fmt.Println("[dpController] Run")
	go func() {
		defer cancel()
		for {
			select {
			case <-ctx.Done():
				return
			default:
				if dc.queue.Len() == 0 {
					//time.Sleep(3 * time.Second)
					continue
				}
				obj, ok := dc.queue.Get()
				if ok { //此时队列为空
					time.Sleep(3 * time.Second)
				}
				dp := obj.(*config.Deployment)
				dc.Sync(dp)
			}
		}
	}()
}

func (dc *DeploymentController) ListDeployments() []*config.Deployment {
	fmt.Println("[dpController] ListDeployments")
	result := make([]*config.Deployment, 0)
	for _, dp := range dc.replicaMap {
		result = append(result, dp)
	}
	return result
}

func (dc *DeploymentController) Sync(dp *config.Deployment) {
	fmt.Println("[dpController] Sync")
	pdw, pdwo := dc.GetPodsWithOwnership(dp)
	runningReplicas := len(pdw)
	dp.Status.Replicas = int32(runningReplicas)
	if dp.Status.Replicas != dp.Spec.Replicas {
		if dp.Status.Replicas < dp.Spec.Replicas {
			dc.IncreaseReplicaCount(dp, pdwo)
		} else if dp.Status.Replicas > dp.Spec.Replicas {
			dc.DecreaseReplicaCount(dp, pdw)
		}
		dp.Status.Replicas = dp.Spec.Replicas
		url := dc.deployClient.BuildURL(apiClient.Create)
		buf, err := dp.JsonMarshal()
		if err != nil {
			fmt.Println(err)
		}
		dc.deployClient.Put(url, buf)
	}
}

func (dc *DeploymentController) SelectDpByLabelSelector(labelSelector selector.LabelSelector) []*config.Deployment {
	fmt.Println("[dpController] SelectDpByLabelSelector")
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
	fmt.Println("[dpController] GetDpsByPod")
	dps := dc.replicaInformer.List()
	fmt.Println("[dpController] GetDpsByPod debugging : Pod's label is ", pod.Metadata.Labels)
	var result []*config.Deployment
	for _, dp := range dps {
		actualDp := dp.(*config.Deployment)
		fmt.Println("[dpController] GetDpsByPod debugging : ", "dp's name is ", actualDp.GetName(), "dp's label is ", actualDp.Spec.Template.Metadata.Labels)
		if selector.LabelCompare(actualDp.Spec.Template.Metadata.Labels, pod.Metadata.Labels) {
			result = append(result, actualDp)
			fmt.Println("[dpController] GetDpsByPod debugging : adding dp", "dp's name is ", actualDp.GetName())
		}
	}
	return result
}

func (dc *DeploymentController) GetPodsWithOwnership(dp *config.Deployment) ([]*config.Pod, []*config.Pod) {
	fmt.Println("[dpController] GetPodsWithOwnership")
	pods := dc.podInformer.List()
	podsWithOwnership := make([]*config.Pod, 0)
	podsWithoutOwnership := make([]*config.Pod, 0)
	//podsPreOwned = make([]core.Pod, 0)
	for _, pod := range pods {
		p := pod.(*config.Pod)
		if selector.LabelSelectorCompare(dp.Spec.Selector, p.Metadata.Labels) {
			if IsDPOwned(dp, p) {
				podsWithOwnership = append(podsWithOwnership, p)
			} else {
				podsWithoutOwnership = append(podsWithoutOwnership, p)
			}
		}
	}
	return podsWithOwnership, podsWithoutOwnership
}

func IsDPOwned(dp *config.Deployment, pd *config.Pod) bool {
	refs := pd.Metadata.OwnerReferences
	if refs == nil || len(refs) == 0 {
		return false
	}
	if refs[0].Kind != string(types.DeploymentObjectType) || refs[0].UID != dp.GetUID() {
		return false
	}
	return true
}

func (dc *DeploymentController) IncreaseReplicaCount(dp *config.Deployment, pdwo []*config.Pod) {
	fmt.Println("[dpController] IncreaseReplicaCount")
	replica := dp.Status.Replicas
	target := dp.Spec.Replicas
	delta := target - replica
	if delta < int32(len(pdwo)) {
		//此时从pdwo中增加即可
		owenerRef := meta.OwnerReference{
			Name:       "",
			UID:        dp.GetUID(),
			APIGroup:   "",
			Kind:       string(types.DeploymentObjectType),
			Controller: false,
		}
		refs := make([]meta.OwnerReference, 0)
		refs = append(refs, owenerRef)
		for i := 0; i < int(delta); i++ {
			pod := pdwo[i]
			pod.Metadata.OwnerReferences = refs
			url := dc.podClient.BuildURL(apiClient.Create)
			buf, err := pod.JsonMarshal()
			if err != nil {
				fmt.Println(err)
			}
			dc.podClient.Put(url, buf)
			replica++
		}

	} else {
		owenerRef := meta.OwnerReference{
			Name:       "",
			UID:        dp.GetUID(),
			APIGroup:   "",
			Kind:       string(types.DeploymentObjectType),
			Controller: false,
		}
		refs := make([]meta.OwnerReference, 0)
		refs = append(refs, owenerRef)
		for _, pod := range pdwo {
			pod.Metadata.OwnerReferences = refs
			url := dc.podClient.BuildURL(apiClient.Create)
			buf, err := pod.JsonMarshal()
			if err != nil {
				fmt.Println(err)
			}
			dc.podClient.Put(url, buf)
			replica++
		}
		mt := dp.Spec.Template.Metadata
		pod := config.Pod{
			ApiVersion: "v1",
			Kind:       string(types.PodObjectType),
			Metadata:   dp.Spec.Template.Metadata,
			Spec:       dp.Spec.Template.Spec,
			Status:     status.PodStatus{},
		}
		for i := 0; i < int(delta)-len(pdwo); i++ {
			mt.Name = "deployment-" + pod.Spec.Containers[0].Name + "-" + testutil.GenerateRandomAlphaOnlyString(5)
			pod.Metadata = mt
			url := dc.podClient.BuildURL(apiClient.Create)
			buf, err := pod.JsonMarshal()
			if err != nil {
				fmt.Println(err)
			}
			dc.podClient.Put(url, buf)
		}
	}
	time.Sleep(1 * time.Second)
}

func (dc *DeploymentController) DecreaseReplicaCount(dp *config.Deployment, pdw []*config.Pod) {
	fmt.Println("[dpController] DecreaseReplicaCount")
	replica := dp.Status.Replicas
	target := dp.Spec.Replicas
	delta := replica - target
	for i := 0; i < int(delta); i++ {
		url := dc.podClient.BuildURL(apiClient.Delete)
		dc.podClient.Delete(url, nil)
	}

}
