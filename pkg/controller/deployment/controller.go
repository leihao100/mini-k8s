package deployment

import (
	"MiniK8S/pkg/api/config"
	"MiniK8S/pkg/api/meta"
	"MiniK8S/pkg/api/status"
	"MiniK8S/pkg/api/types"
	"MiniK8S/pkg/apiClient"
	"context"
	"github.com/google/uuid"
	"strconv"
)

type DeploymentController struct {
	podClient    *apiClient.Client
	deployClient *apiClient.Client
	replicaMap   map[uuid.UUID]*config.Deployment
	pods         map[uuid.UUID][]*config.Pod
}

func NewController() *DeploymentController {
	dpClient := apiClient.NewRESTClient(types.DeploymentObjectType)
	pdClient := apiClient.NewRESTClient(types.PodObjectType)
	return &DeploymentController{
		deployClient: dpClient,
		podClient:    pdClient,
		replicaMap:   make(map[uuid.UUID]*config.Deployment),
	}
}

func (dc *DeploymentController) Run(ctx context.Context) {

}

func (dc *DeploymentController) GetDeployment(uid uuid.UUID) *config.Deployment {
	return dc.replicaMap[uid]
}

func (dc *DeploymentController) ListDeployments() []*config.Deployment {
	result := make([]*config.Deployment, 0)
	for _, dp := range dc.replicaMap {
		result = append(result, dp)
	}
	return result
}

func (dc *DeploymentController) CreateDeployment(dp *config.Deployment) {
	dc.replicaMap[dp.Metadata.Uid] = dp
	replicas := dp.Spec.Replicas
	owner := make([]meta.OwnerReference, 0)
	owner = append(owner, meta.OwnerReference{
		Name:       dp.Metadata.Name,
		UID:        dp.Metadata.Uid,
		Kind:       "Deployment",
		Controller: true,
		APIGroup:   dp.ApiVersion,
	})
	for i := 0; i < int(replicas); i++ {
		URL := dc.podClient.BuildURL(apiClient.Create)
		spec := dp.Spec.Template.Spec
		pod := config.Pod{
			ApiVersion: "",
			Kind:       "pod",
			Metadata: meta.ObjectMeta{
				Name:            dp.Spec.Template.Metadata.Name + "_" + strconv.Itoa(i),
				Namespace:       dp.Spec.Template.Metadata.Namespace,
				Labels:          dp.Spec.Template.Metadata.Labels,
				OwnerReferences: owner,
			},
			Spec:   spec,
			Status: status.PodStatus{},
		}
		buf := pod.Marshal()
		dc.podClient.Post(URL, buf)
		dc.pods[dp.Metadata.Uid] = append(dc.pods[dp.Metadata.Uid], &pod)
	}
}

func (dc *DeploymentController) UpdateDeployment(dp, new_dp *config.Deployment) {
	//var isNum=false
	//var isTem=false
	//if dp.Spec.Replicas!=new_dp.Spec.Replicas {
	//	isNum=true
	//}
	//if new_dp.Spec.Template!=dp.Spec.Template

}
func (dc *DeploymentController) AddPod(pod *config.Pod) {
	if pod.Metadata.OwnerReferences[0].Controller == false {
		return
	}
	if pod.Metadata.OwnerReferences[0].Kind != "Deployment" {
		return
	}
	uid := pod.Metadata.OwnerReferences[0].UID
	dp := dc.GetDeployment(uid)
	if dp == nil {
		//如果该pod对应的DP不存在，则删除该pod
		poduid := pod.Metadata.Uid
		url := dc.podClient.BuildURL(apiClient.Delete) + "/" + poduid.String()
		dc.podClient.Delete(url, nil)
		return
	}
	dc.pods[uid] = append(dc.pods[uid], pod)
}

func (dc *DeploymentController) DeletePod(pod *config.Pod) {
	if pod.Metadata.OwnerReferences[0].Controller == false {
		return
	}
	if pod.Metadata.OwnerReferences[0].Kind != "Deployment" {
		return
	}
	uid := pod.Metadata.OwnerReferences[0].UID
	dp := dc.GetDeployment(uid)
	if dp == nil {
		return
	}
	for i, poddp := range dc.pods[uid] {
		if poddp.Metadata.Uid == pod.Metadata.Uid {
			dc.pods[uid] = append(dc.pods[uid][:i], dc.pods[uid][i+1:]...)
			return
		}
	}
}

func (dc *DeploymentController) DeleteDeployment(uid uuid.UUID) {
	for _, pod := range dc.pods[uid] {
		poduid := pod.Metadata.Uid
		url := dc.podClient.BuildURL(apiClient.Delete) + "/" + poduid.String()
		dc.podClient.Delete(url, nil)
	}
	delete(dc.pods, uid)
	delete(dc.replicaMap, uid)

}

func (dc *DeploymentController) Scale(uid uuid.UUID) {
	//num := dc.replicaMap[uid].Spec.Replicas

}

func (dc *DeploymentController) Sync() {

}

func (dc *DeploymentController) Watch() {

}
