package hpa

import (
	"MiniK8S/pkg/api/config"
	"MiniK8S/pkg/api/selector"
	"MiniK8S/pkg/api/types"
	"MiniK8S/pkg/apiClient"
	"MiniK8S/pkg/controller/cache"
	"MiniK8S/pkg/controller/hpa/metrics"
	"context"
	"fmt"
	"math"
	"slices"
	"strings"
	"time"
)

var DefaultScaleUpPolicy = config.HPAScalingRules{
	Policies: []config.HPAScalingPolicy{{
		Type:          config.PolicyPercent,
		Value:         100,
		PeriodSeconds: 15,
	}},
	SelectPolicy:               config.Max,
	StabilizationWindowSeconds: 0,
}

var DefaultScaleDownPolicy = config.HPAScalingRules{
	Policies: []config.HPAScalingPolicy{{
		Type:          config.PolicyPercent,
		Value:         33,
		PeriodSeconds: 15,
	}},
	SelectPolicy:               config.Max,
	StabilizationWindowSeconds: 300,
}

type HpaController struct {
	podClient    *apiClient.Client
	hpaClient    *apiClient.Client
	deployClient *apiClient.Client
	podInformer  *cache.Informer
	hpaInformer  *cache.Informer
	dpInformer   *cache.Informer
	queue        *cache.WorkQueue
	metricClient *metrics.HPAMetricsClient
}

func NewController(pi *cache.Informer, hi *cache.Informer, di *cache.Informer, pc *apiClient.Client, hc *apiClient.Client, dc *apiClient.Client, nc *apiClient.Client) *HpaController {
	hpc := &HpaController{
		hpaClient:    hc,
		podClient:    pc,
		deployClient: dc,
		podInformer:  pi,
		hpaInformer:  hi,
		dpInformer:   di,
		metricClient: metrics.NewMetricsClient(nc),
		queue:        cache.NewWorkQueue(),
	}
	hpc.hpaInformer.AddEventHandler(cache.EventHandlerFuncs{
		AddFunc:    hpc.AddHpa,
		UpdateFunc: hpc.UpdateHpa,
		DeleteFunc: hpc.DeleteHpa,
	})
	return hpc
}

func (hpc *HpaController) AddHpa(obj interface{}) {
	hpa := obj.(*config.HorizontalPodAutoscaler)
	hpc.queue.Add(hpa)
}

func (hpc *HpaController) UpdateHpa(oldObj, newObj interface{}) {
	hpa := oldObj.(*config.HorizontalPodAutoscaler)
	hpc.queue.Add(hpa)
}

func (hpc *HpaController) DeleteHpa(obj interface{}) {
	//hpa := obj.(*config.HorizontalPodAutoscaler)
	//handle delete hpa

}

func (hpc *HpaController) CalculateTarget(hpa *config.HorizontalPodAutoscaler, realnum int, metric metrics.PodMetric) int {
	fmt.Println("[hpaController] CalculateTarget")
	res := make([]int, 0)

	for _, met := range hpa.Spec.Metrics {
		res = append(res, hpc.CalculateTargetByOneType(realnum, metric, met))
	}
	desire := int32(slices.Max(res))
	if desire < hpa.Spec.MinReplicas {
		desire = hpa.Spec.MinReplicas
	}
	if desire > hpa.Spec.MaxReplicas {
		desire = hpa.Spec.MaxReplicas
	}
	return int(desire)

}

func (hpc *HpaController) Run(ctx context.Context, cancel context.CancelFunc) {
	fmt.Println("[hpa] Starting HpaController")

	go func() {
		defer cancel()
		for {
			select {
			case <-ctx.Done():
				fmt.Println("[hpa] Stopping HpaController")
				return
			default:
				//if hpc.queue.Len() == 0 {
				//	time.Sleep(3 * time.Second)
				//}
				obj, _ := hpc.queue.Get()
				//if !ok {
				//	time.Sleep(3 * time.Second)
				//}
				hpa := obj.(*config.HorizontalPodAutoscaler)
				hpc.Sync(hpa)
				//time.Sleep(1000 * time.Millisecond)
			}
		}
	}()
}

func (hpc *HpaController) CalculateTargetByOneType(realnum int, metric metrics.PodMetric, spec config.MetricSpec) int {
	fmt.Println("[hpaController] CalculateTargetByOneType")
	ty := metrics.ResourceType(spec.Type)
	target := spec.Target
	var res = 0

	switch ty {
	case metrics.CPU:
		res = int(math.Ceil(float64(realnum) * (float64(metric.CPU) / float64(target))))
	case metrics.Memory:
		res = int(math.Ceil(float64(realnum) * (float64(metric.Memory) / float64(target))))
	default:
		fmt.Println("Unsupported metric type:", ty)
	}

	return res
}

func (hpc *HpaController) Sync(hpa *config.HorizontalPodAutoscaler) {
	fmt.Println("[hpaController] Syncing Hpa")
	hpc.metricClient.Sync()
	ContainerInfo := hpc.metricClient.GetContainerMetricMap()
	metric, i, err := hpc.GetResourceMetric(hpa, ContainerInfo)
	if err != nil {
		panic(err)
		return
	}
	desire := int32(hpc.CalculateTarget(hpa, i, metric))
	if desire > (hpa.Spec.MaxReplicas) {
		desire = (hpa.Spec.MaxReplicas)
	}
	if desire < (hpa.Spec.MinReplicas) {
		desire = (hpa.Spec.MinReplicas)
	}
	replicas, err := hpc.GetHpaReplicas(hpa)
	if err != nil {
		return
	}
	if int32(desire) == replicas {
		return
	}
	var rule config.HPAScalingRules
	if int32(desire) < replicas {
		//芝士缩容
		rule = hpa.Spec.Behavior.ScaleDown
		if rule.SelectPolicy == "" {
			rule = DefaultScaleDownPolicy
		}
		if int32(time.Since(hpa.Status.LastScaleTime).Seconds()) < rule.StabilizationWindowSeconds {
			//未到冷却时间
			return
		}
		for _, policy := range rule.Policies {
			switch policy.Type {
			case config.PolicyPod:
				if replicas-desire > policy.Value {
					desire = replicas - policy.Value
				}
			case config.PolicyPercent:
				if float64(replicas-desire)/float64(replicas)*100 > float64(policy.Value) {
					desire = replicas - int32(float64(policy.Value)/100*float64(replicas))
				}
			}

		}
	} else {
		//芝士扩容
		rule = hpa.Spec.Behavior.ScaleUp
		if rule.SelectPolicy == "" {
			rule = DefaultScaleUpPolicy
		}
		if int32(time.Since(hpa.Status.LastScaleTime).Seconds()) < rule.StabilizationWindowSeconds {
			//未到冷却时间
			return
		}
		for _, policy := range rule.Policies {
			switch policy.Type {
			case config.PolicyPod:
				if desire-replicas > policy.Value {
					desire = replicas + policy.Value
				}
			case config.PolicyPercent:
				if float64(desire-replicas)/float64(replicas)*100 > float64(policy.Value) {
					desire = replicas + int32(float64(policy.Value)/100*float64(replicas))
				}
			}
		}
	}

	hpc.Scale(hpa, int(desire))
	hpa.Status.CurrentReplicas = int32(desire)
	url := hpc.hpaClient.BuildURL(apiClient.Create)
	buf, _ := hpa.JsonMarshal()
	hpc.hpaClient.Put(url, buf)
	//choose behaviour

}

func (hpc *HpaController) Scale(hpa *config.HorizontalPodAutoscaler, desire int) {

	dpName := hpa.Spec.ScaleTargetRef.Name
	dps := hpc.dpInformer.List()

	for _, dp := range dps {
		d := dp.(*config.Deployment)
		if d.Metadata.Name == dpName {
			//scale

			d.Spec.Replicas = int32(desire)
			bytes, err := d.JsonMarshal()
			if err != nil {
				fmt.Println("[hpa] Failed to scale deployment because of marshal:", err)
			}
			url := hpc.deployClient.BuildURL(apiClient.Create)
			hpc.deployClient.Put(url, bytes)
			return
		}
	}
}

func (hpc *HpaController) GetResourceMetric(hpa *config.HorizontalPodAutoscaler, cmap metrics.ConatinerMetricsInfo) (metrics.PodMetric, int, error) {
	if !strings.EqualFold(hpa.Spec.ScaleTargetRef.Kind, string(types.DeploymentObjectType)) {
		return metrics.PodMetric{}, 0, fmt.Errorf("HPA is not a deployment object")
	}
	dpName := hpa.Spec.ScaleTargetRef.Name
	dps := hpc.dpInformer.List()

	for _, dp := range dps {
		d := dp.(*config.Deployment)
		if d.Metadata.Name == dpName {
			return hpc.GetDeploymentResourceMetric(d, cmap)
		}
	}

	return metrics.PodMetric{}, 0, fmt.Errorf("do not find deployment belongs to this HPA")
}

func (hpc *HpaController) GetDeploymentResourceMetric(dp *config.Deployment, cmap metrics.ConatinerMetricsInfo) (metrics.PodMetric, int, error) {
	fmt.Println("[hpaController] GetDeploymentResourceMetric")
	pods := hpc.GetPodFromDeployment(dp)
	pdNum := len(pods)
	res := make([]metrics.PodMetric, 0)
	var metricsVector []metrics.PodMetric

	for _, pod := range pods {
		metricsVector = make([]metrics.PodMetric, 0)
		for _, container := range pod.Spec.Containers {
			metricsVector = append(metricsVector, cmap[container.Name])
		}
		res = append(res, metrics.PodMetricsSum(metricsVector))
	}

	return metrics.CalculateAverage(res), pdNum, nil
}

func (hpc *HpaController) GetPodFromDeployment(dp *config.Deployment) []config.Pod {
	fmt.Println("[hpaController] GetPodFromDeployment")
	pods := hpc.podInformer.List()
	res := make([]config.Pod, 0)

	for _, pod := range pods {
		p := pod.(*config.Pod)
		if selector.LabelSelectorCompare(dp.Spec.Selector, p.Metadata.Labels) {
			res = append(res, *p)
		}
	}

	return res
}

func (hpc *HpaController) GetHpaReplicas(hpa *config.HorizontalPodAutoscaler) (int32, error) {
	dpl := hpc.dpInformer.List()
	res := int32(-1)
	for _, d := range dpl {
		dp := d.(*config.Deployment)
		if dp.Metadata.Name == hpa.Spec.ScaleTargetRef.Name {
			res = dp.Spec.Replicas
		}
	}
	if res == -1 {
		return res, fmt.Errorf("[hpa Controller] Failed to find a deployment belongs to this HPA")
	}
	return res, nil
}
