package metrics

import (
	"MiniK8S/pkg/api/config"
	"MiniK8S/pkg/api/selector"
	apitypes "MiniK8S/pkg/api/types"
	"MiniK8S/pkg/apiClient"
	"MiniK8S/pkg/apiClient/listwatch"
	"MiniK8S/pkg/kubelet/cadvisor"
	"context"
	"fmt"
	"github.com/google/cadvisor/info/v1"
	"time"
)

const CAdvisorPort = ":9090"

type HPAMetricsClient struct {
	nodeClient      *apiClient.Client
	nodeListWatcher listwatch.ListerWatcher
	//将informer作为参数传入太麻烦了，使用 lw调用list即可
	//nodeInformer   	*cache.Informer
	metricsInfo ConatinerMetricsInfo
	//we use node name to detect cadvisor
	cadvisorClients map[string]*cadvisor.CAdvisorClient
}

func NewMetricsClient(nodeClient *apiClient.Client) *HPAMetricsClient {
	return &HPAMetricsClient{
		metricsInfo: make(ConatinerMetricsInfo),
	}
}
func (hmc *HPAMetricsClient) GetResourceMetric(ctx context.Context, resource apitypes.ApiObjectType, namespace string, selector selector.LabelSelector, container string) (ConatinerMetricsInfo, time.Time, error) {
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

func (hmc *HPAMetricsClient) GetPodMetrics(pod *config.Pod) []PodMetric {
	nodename := pod.Spec.NodeName
	cadvisorClient, ok := hmc.cadvisorClients[nodename]
	if !ok {
		fmt.Println("[metrics][GetPodMetric] cadvisorClient not found")
		return []PodMetric{}
	}
	query := v1.ContainerInfoRequest{
		NumStats: 12,
		Start:    time.Time{},
		End:      time.Time{},
	}
	metricsInfo := make([]PodMetric, 0)
	infos, err := cadvisorClient.Inspect(&query)
	for _, info := range infos {
		for _, status := range pod.Status.ContainerStatuses {
			if status.Name == info.Aliases[0] {
				metricsInfo = append(metricsInfo, PodMetric{
					Timestamp: info.Stats[0].Timestamp,
					Window:    0,
					CPU:       info.Stats[0].Cpu.Usage.Total,
					Memory:    info.Stats[0].Memory.Usage,
				})
			}
		}
	}

	if err != nil {
		return []PodMetric{}
	}
	return metricsInfo
}

func PodMetricsSum(pms []PodMetric) PodMetric {
	result := PodMetric{
		Timestamp: time.Now(),
		Window:    0,
		CPU:       0,
		Memory:    0,
	}
	for _, pm := range pms {
		result.CPU += pm.CPU
		result.Memory += pm.Memory
	}
	return result
}

func AddPodMetric(a, b PodMetric) PodMetric {
	return PodMetric{
		Timestamp: time.Time{},
		Window:    0,
		CPU:       a.CPU + b.CPU,
		Memory:    a.Memory + b.Memory,
	}
}

func DividePodMetric(a PodMetric, divider uint64) PodMetric {
	a.CPU /= divider
	a.Memory /= divider
	return a
}

func CalculateAverage(pms []PodMetric) PodMetric {
	var cnt uint64 = 0
	result := PodMetric{
		Timestamp: time.Now(),
		Window:    0,
		CPU:       0,
		Memory:    0,
	}
	for _, pm := range pms {
		result.CPU += pm.CPU
		result.Memory += pm.Memory
		cnt++
	}
	result.CPU /= cnt
	result.Memory /= cnt
	return result
}

func (hmc *HPAMetricsClient) Sync() {

	hmc.cadvisorClients = make(map[string]*cadvisor.CAdvisorClient)
	nodes, _ := hmc.nodeListWatcher.List(config.ListOptions{
		Watch: false,
	})
	ns := nodes.GetItems()
	for _, n := range ns {
		node := n.(*config.Node)
		url := "http://" + node.Status.Addresses.Address + CAdvisorPort
		hmc.cadvisorClients[node.Metadata.Name] = cadvisor.NewCAdvisor(url)
	}
	for _, client := range hmc.cadvisorClients {
		query := v1.ContainerInfoRequest{
			NumStats: 12,
			Start:    time.Time{},
			End:      time.Time{},
		}
		infos, err := client.Inspect(&query)
		if err != nil {
			fmt.Println("[metrics client][Sync]", err)
		}
		for _, info := range infos {
			hmc.metricsInfo[info.Aliases[0]] = PodMetric{
				Timestamp: info.Stats[0].Timestamp,
				Window:    0,
				CPU:       info.Stats[0].Cpu.Usage.Total,
				Memory:    info.Stats[0].Memory.Usage,
			}
		}
	}
}

func (hmc *HPAMetricsClient) GetContainerMetricMap() ConatinerMetricsInfo {
	return hmc.metricsInfo
}
