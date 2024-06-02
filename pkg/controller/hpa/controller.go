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
	"time"
)

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

func NewController(pi *cache.Informer, hi *cache.Informer, pc *apiClient.Client, hc *apiClient.Client, dc *apiClient.Client) *HpaController {
	hpc := &HpaController{
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
	res := make([]int, 0)

	for _, met := range hpa.Spec.Metrics {
		res = append(res, hpc.CalculateTargetByOneType(realnum, metric, met))
	}
	return slices.Max(res)

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
				obj, ok := hpc.queue.Get()
				if !ok {
					time.Sleep(3 * time.Second)
				}
				hpa := obj.(*config.HorizontalPodAutoscaler)
				hpc.Sync(hpa)
			}
		}
	}()
}

func (hpc *HpaController) CalculateTargetByOneType(realnum int, metric metrics.PodMetric, spec config.MetricSpec) int {
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
	hpc.metricClient.Sync()
	ContainerInfo := hpc.metricClient.GetContainerMetricMap()
	metric, i, err := hpc.GetResourceMetric(hpa, ContainerInfo)
	if err != nil {
		panic(err)
		return
	}
	desire := hpc.CalculateTarget(hpa, i, metric)
	if desire > int(hpa.Spec.MaxReplicas) {
		desire = int(hpa.Spec.MaxReplicas)
	}
	if desire < int(hpa.Spec.MinReplicas) {
		desire = int(hpa.Spec.MinReplicas)
	}
	hpc.Scale(hpa, desire)
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
			//todo update scale time and add scale methods
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
	if hpa.Spec.ScaleTargetRef.Kind != string(types.DeploymentObjectType) {
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
