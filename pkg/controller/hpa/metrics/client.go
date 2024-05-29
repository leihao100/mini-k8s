package metrics

import (
	"MiniK8S/pkg/api/config"
	"MiniK8S/pkg/api/selector"
	apitypes "MiniK8S/pkg/api/types"
	"MiniK8S/pkg/apiClient"
	"MiniK8S/pkg/kubelet/cadvisor"
	"context"
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
