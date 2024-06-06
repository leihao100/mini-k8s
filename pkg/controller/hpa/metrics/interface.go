package metrics

import (
	"MiniK8S/pkg/api/config"
	"MiniK8S/pkg/api/selector"
	apitypes "MiniK8S/pkg/api/types"
	"context"
	"time"
)

type ResourceType string

const (
	CPU    ResourceType = "cpu"
	Memory ResourceType = "memory"
)

// PodMetric contains pod metric value (the metric values are expected to be the metric as a milli-value)
type PodMetric struct {
	Timestamp time.Time
	Window    time.Duration
	CPU       uint64
	Memory    uint64
}

// ConatinerMetricsInfo contains pod metrics as a map from Container name to PodMetric
type ConatinerMetricsInfo map[string]PodMetric

// MetricsClient knows how to query a remote interface to retrieve container-level
// resource metrics as well as pod-level arbitrary metrics
type MetricsClient interface {
	// GetResourceMetric gets the given resource metric (and an associated oldest timestamp)
	// for the specified named container in all pods matching the specified selector in the given namespace and when
	// the container is an empty string it returns the sum of all the container metrics.
	GetResourceMetric(ctx context.Context, resource apitypes.ApiObjectType, namespace string, selector selector.LabelSelector, container string) (ConatinerMetricsInfo, time.Time, error)

	// GetRawMetric gets the given metric (and an associated oldest timestamp)
	// for all pods matching the specified selector in the given namespace
	GetRawMetric(metricName string, namespace string, selector selector.LabelSelector, metricSelector selector.LabelSelector) (ConatinerMetricsInfo, time.Time, error)

	// GetObjectMetric gets the given metric (and an associated timestamp) for the given
	// object in the given namespace
	GetObjectMetric(metricName string, namespace string, objectRef config.HorizontalPodAutoscaler, metricSelector selector.LabelSelector) (int64, time.Time, error)

	// GetExternalMetric gets all the values of a given external metric
	// that match the specified selector.
	GetExternalMetric(metricName string, namespace string, selector selector.LabelSelector) ([]int64, time.Time, error)
}
