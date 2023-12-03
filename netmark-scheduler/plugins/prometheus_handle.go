package plugins

import (
	"fmt"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"k8s.io/klog"
	"strconv"
	"strings"
	"time"

	"context"
)

const (
	nodeMeasureQuery      = "sum_over_time(process_cpu_seconds_total[1y])"
	IstioRequestTotal     = "istio_requests_total{pod=\"%s\"}"
	IstioRequestBytesSum  = "istio_request_bytes_sum{pod=\"%s\",source_app=\"%s\",destination_app=\"%s\"}"
	IstioResponseBytesSum = "istio_response_bytes_sum{pod=\"%s\",source_app=\"%s\",destination_app=\"%s\"}"
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

func ExtractFields(input string) (string, string, bool) {
	// Split the string into key-value pairs
	pairs := strings.Split(input, ", ")
	var destinationApp, sourceApp string

	for _, pair := range pairs {
		kv := strings.Split(pair, "=")
		key := kv[0]
		value := strings.Trim(kv[1], "\"")

		// Check for the response_code
		if key == "response_code" && value != "200" {
			return "", "", false
		}

		if value != "unknown" {
			if key == "destination_app" {
				destinationApp = value
			} else if key == "source_app" {
				sourceApp = value
			}
		}
	}

	return destinationApp, sourceApp, true
}

func hasSuccessResponseCode(input string) bool {
	// Split the string into key-value pairs
	pairs := strings.Split(input, ", ")
	for _, pair := range pairs {
		kv := strings.Split(pair, "=")
		key := kv[0]
		value := strings.Trim(kv[1], "\"")

		// Check for the response_code
		if key == "response_code" && value == "200" {
			return true
		}
	}

	return false
}

func extractAppName(podName string) string {
	pairs := strings.Split(podName, "-")
	if len(pairs) < 3 {
		return ""
	}

	return pairs[0] + "-" + pairs[1]
}

func (p *PrometheusHandle) GetIstioMetric(isRequest bool, sourcePod, destinationPod string) float64 {
	var query string
	if isRequest {
		query = IstioRequestBytesSum
	} else {
		query = IstioResponseBytesSum
	}

	sourceApp, destinationApp := extractAppName(sourcePod), extractAppName(destinationPod)

	value, err := p.query(fmt.Sprintf(query, sourcePod, sourceApp, destinationApp))
	if err != nil {
		fmt.Printf("[NetworkTraffic Plugin] Error querying prometheus: %s\n", err.Error())
		return 0
	}

	nodeMeasure := value.(model.Vector)
	for _, res := range nodeMeasure {
		ok := hasSuccessResponseCode(res.String())

		if ok {
			if s, err := strconv.ParseFloat(res.Value.String(), 64); err == nil {
				return s
			}
		}
	}

	return 0
}

func (p *PrometheusHandle) GetNeighborServices(podName string) []string {
	value, err := p.query(fmt.Sprintf(IstioRequestTotal, podName))
	if err != nil {
		fmt.Printf("[NetworkTraffic Plugin] Error querying prometheus: %s\n", err.Error())
		return nil
	}

	nodeMeasure := value.(model.Vector)
	if len(nodeMeasure) == 0 {
		fmt.Printf("[NetworkTraffic Plugin] Empty response. Response Len: %d", len(nodeMeasure))
		return nil
	}

	var queryResult []string
	for _, res := range nodeMeasure {
		destinationApp, sourceApp, ok := ExtractFields(res.String())
		if ok && destinationApp != "" && !strings.HasPrefix(podName, destinationApp) {
			queryResult = append(queryResult, destinationApp)
		}

		if ok && sourceApp != "" && !strings.HasPrefix(podName, sourceApp) {
			queryResult = append(queryResult, sourceApp)
		}
	}

	return queryResult
}

// isNetmarksPod checks if the pod's labels match the 'app:netmarks-<serviceName>' pattern
func isNetmarksPod(label string) bool {
	if strings.HasPrefix(label, "netmarks-") {
		return true
	}

	return false
}

// CalculateTrafficScore calculates traffic flow over the pod's entire lifetime
// TODO: add time period
func (p *PrometheusHandle) CalculateTrafficScore(pod1 string, pod2 string) float64 {
	requestByteSum := p.GetIstioMetric(true, pod1, pod2)
	fmt.Printf("Request from %s to %s: %s\n", pod1, pod2, strconv.FormatFloat(requestByteSum, 'f', -1, 64))

	responseByteSum := p.GetIstioMetric(false, pod2, pod1)
	fmt.Printf("Response from %s to %s: %s\n", pod2, pod1, strconv.FormatFloat(responseByteSum, 'f', -1, 64))

	return requestByteSum + responseByteSum
}

func (p *PrometheusHandle) GetGauge(node string) (*model.Sample, error) {
	value, err := p.query(fmt.Sprintf(nodeMeasureQuery))

	//fmt.Println(fmt.Sprintf(nodeMeasureQuery, p.deviceName, p.timeRange, node))

	klog.Infof("[NetworkTraffic] Prometheus Query: %s", nodeMeasureQuery)
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
