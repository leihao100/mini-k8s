package metrics

import (
	"fmt"
	"github.com/google/uuid"
)

// GetResourceUtilizationRatio takes in a set of metrics, a set of matching requests,
// and a target utilization percentage, and calculates the ratio of
// desired to actual utilization (returning that, the actual utilization, and the raw average value)
func GetResourceUtilizationRatio(metrics PodMetricsInfo, ty ResourceType, requests map[uuid.UUID]int64, targetUtilization int32) (utilizationRatio float64, currentUtilization int32, rawAverageValue int64, err error) {
	metricsTotal := int64(0)
	requestsTotal := int64(0)
	numEntries := 0

	for podName, metric := range metrics {
		request, hasRequest := requests[podName]
		if !hasRequest {
			// we check for missing requests elsewhere, so assuming missing requests == extraneous metrics
			continue
		}
		switch ty {
		case CPU:
			metricsTotal += metric.CPU
		case Memory:
			metricsTotal += metric.Memory
		}
		//metricsTotal += metric.Value
		requestsTotal += request
		numEntries++
	}

	// if the set of requests is completely disjoint from the set of metrics,
	// then we could have an issue where the requests total is zero
	if requestsTotal == 0 {
		return 0, 0, 0, fmt.Errorf("no metrics returned matched known pods")
	}

	currentUtilization = int32((metricsTotal * 100) / requestsTotal)

	return float64(currentUtilization) / float64(targetUtilization), currentUtilization, metricsTotal / int64(numEntries), nil
}

// GetMetricUsageRatio takes in a set of metrics and a target usage value,
// and calculates the ratio of desired to actual usage
// (returning that and the actual usage)
func GetMetricUsageRatio(metrics PodMetricsInfo, ty ResourceType, targetUsage int64) (usageRatio float64, currentUsage int64) {
	metricsTotal := int64(0)
	for _, metric := range metrics {
		switch ty {
		case CPU:
			metricsTotal += metric.CPU
		case Memory:
			metricsTotal += metric.Memory
		}
	}

	currentUsage = metricsTotal / int64(len(metrics))

	return float64(currentUsage) / float64(targetUsage), currentUsage
}
