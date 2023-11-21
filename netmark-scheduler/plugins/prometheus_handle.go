package plugins

import (
	"fmt"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"k8s.io/klog"
	"time"

	"context"
)

const (
	// nodeMeasureQueryTemplate is the template string to get the query for the node used bandwidth
	// nodeMeasureQueryTemplate = "sum_over_time(node_network_receive_bytes_total{device=\"%s\"}[%ss])"
	//nodeMeasureQueryTemplate = "sum_over_time(node_network_receive_bytes_total{device=\"%s\"}[%ss]) * on(instance) group_left(nodename) (node_uname_info{nodename=\"%s\"})"
	//nodeMeasureQueryTemplate = "istio_request_bytes_sum"
	//nodeMesaureQuery = "istio_request_bytes_sum"
	nodeMesaureQuery = "sum_over_time(process_cpu_seconds_total[1y])"
)

type PrometheusHandle struct {
	deviceName string
	timeRange  time.Duration
	ip         string
	client     v1.API
}

func NewProme(ip, deviceName string, timeRace time.Duration) *PrometheusHandle {
	client, err := api.NewClient(api.Config{Address: ip})
	if err != nil {
		klog.Fatalf("[NetworkTraffic Plugin] FatalError creating prometheus client: %s", err.Error())
	}
	return &PrometheusHandle{
		deviceName: deviceName,
		ip:         ip,
		timeRange:  timeRace,
		client:     v1.NewAPI(client),
	}
}

func (p *PrometheusHandle) GetGauge(node string) (*model.Sample, error) {
	value, err := p.query(fmt.Sprintf(nodeMesaureQuery))

	//fmt.Println(fmt.Sprintf(nodeMesaureQuery, p.deviceName, p.timeRange, node))

	klog.Infof("[NetworkTraffic] Prometheus Query: %s", nodeMesaureQuery)
	klog.Infof("[NetworkTraffic] Prometheus Device Name: %s", p.deviceName)
	klog.Infof("[NetworkTraffic] Prometheus Time Range: %s", p.timeRange)
	klog.Infof("[NetworkTraffic] Prometheus Node: %s", node)

	if err != nil {
		return nil, fmt.Errorf("[NetworkTraffic Plugin] Error querying prometheus: %w", err)
	}

	nodeMeasure := value.(model.Vector)
	// TODO: What if we have an empty query result?
	if len(nodeMeasure) == 0 {
		return nil, fmt.Errorf("[NetworkTraffic Plugin] Empty response. Response Len: %d", len(nodeMeasure))
	}
	return nodeMeasure[0], nil
}

func (p *PrometheusHandle) query(promQL string) (model.Value, error) {
	results, warnings, err := p.client.Query(context.Background(), promQL, time.Now())
	if len(warnings) > 0 {
		klog.Warningf("[NetworkTraffic Plugin] Warnings: %v\n", warnings)
	}

	return results, err
}
