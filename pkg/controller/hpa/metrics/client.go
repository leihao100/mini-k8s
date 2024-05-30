package metrics

import (
	"MiniK8S/pkg/api/config"
	"MiniK8S/pkg/api/selector"
	apitypes "MiniK8S/pkg/api/types"
	"MiniK8S/pkg/apiClient"
	"MiniK8S/pkg/kubelet/cadvisor"
	"context"
	"github.com/google/cadvisor/info/v1"
	"github.com/google/uuid"
	"time"
)

type HPAMetricsClient struct {
	nodeClient      *apiClient.Client
	metricsInfo     PodMetricsInfo
	cadvisorClients map[uuid.UUID]*cadvisor.CAdvisorClient
}

func NewMetricsClient(nodeClient *apiClient.Client) *HPAMetricsClient {
	return &HPAMetricsClient{
		metricsInfo: make(PodMetricsInfo),
	}
}
func (hmc *HPAMetricsClient) GetResourceMetric(ctx context.Context, resource apitypes.ApiObjectType, namespace string, selector selector.LabelSelector, container string) (PodMetricsInfo, time.Time, error) {
	query := v1.ContainerInfoRequest{
		NumStats: 12,
		Start:    time.Time{},
		End:      time.Time{},
	}
	for _, client := range hmc.cadvisorClients {
		_, err := client.Inspect(&query)
		if err != nil {
			return nil, time.Time{}, err
		}
		//for _, info := range infos {
		//
		//}
	}
	return nil, time.Now(), nil
}
func (hmc *HPAMetricsClient) GetRawMetric(metricName string, namespace string, selector selector.LabelSelector, metricSelector selector.LabelSelector) (PodMetricsInfo, time.Time, error) {
	return nil, time.Now(), nil
}
func (hmc *HPAMetricsClient) GetObjectMetric(metricName string, namespace string, objectRef config.HorizontalPodAutoscaler, metricSelector selector.LabelSelector) (int64, time.Time, error) {
	return 0, time.Now(), nil
}
func (hmc *HPAMetricsClient) GetExternalMetric(metricName string, namespace string, selector selector.LabelSelector) ([]int64, time.Time, error) {
	return nil, time.Now(), nil
}

func AddResource(a, b v1.ContainerInfo) (uint64, uint64) {
	//cpu := a.Stats[0].Cpu.Usage.Total
	return 0, 0
}
