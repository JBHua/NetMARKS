package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/prometheus/client_golang/api"
	promethus_v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/klog"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Scoring struct {
	prometheus *PrometheusHandle
	clientSet  *kubernetes.Clientset
}

type NodeScore struct {
	NodeName string
	PodName  string
	Score    float64
}

const (
	nodeMesaureQuery = "sum_over_time(process_cpu_seconds_total[1y])"
	// IstioResponseByteSum TODO: add time range parameter
	IstioRequestTotal     = "istio_requests_total{pod=\"%s\"}"
	IstioRequestBytesSum  = "istio_request_bytes_sum{pod=\"%s\",source_app=\"%s\",destination_app=\"%s\"}"
	IstioResponseBytesSum = "istio_response_bytes_sum{pod=\"%s\",source_app=\"%s\",destination_app=\"%s\"}"
)

type PrometheusHandle struct {
	client promethus_v1.API
}

func sum(array []float64) float64 {
	var result float64

	for _, v := range array {
		result += v
	}
	return result
}

func NewPrometheusClient(ip string) *PrometheusHandle {
	client, err := api.NewClient(api.Config{Address: ip})
	if err != nil {
		klog.Fatalf("[NetworkTraffic] FatalError creating prometheus client: %s", err.Error())
	}
	return &PrometheusHandle{
		client: promethus_v1.NewAPI(client),
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

func (p *PrometheusHandle) query(promQL string) (model.Value, error) {
	results, warnings, err := p.client.Query(context.Background(), promQL, time.Now())
	if len(warnings) > 0 {
		klog.Warningf("[NetworkTraffic Plugin] Warnings: %v\n", warnings)
	}

	return results, err
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
func (s *Scoring) CalculateTrafficScore(pod1 string, pod2 string) float64 {
	requestByteSum := s.prometheus.GetIstioMetric(true, pod1, pod2)
	fmt.Printf("Request from %s to %s: %s\n", pod1, pod2, strconv.FormatFloat(requestByteSum, 'f', -1, 64))

	responseByteSum := s.prometheus.GetIstioMetric(false, pod2, pod1)
	fmt.Printf("Response from %s to %s: %s\n", pod2, pod1, strconv.FormatFloat(responseByteSum, 'f', -1, 64))

	return requestByteSum + responseByteSum
}

func (s *Scoring) CalculateScore(nodeName string, currentPod string) NodeScore {
	var scores []float64

	podsOnNode, _ := s.clientSet.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{
		FieldSelector: "spec.nodeName=" + nodeName,
	})

	for _, podOnNode := range podsOnNode.Items {
		for _, podTrafficNeighbor := range s.prometheus.GetNeighborServices(currentPod) {
			//only care about traffic neighbor on running on the same node.
			if !strings.HasPrefix(podOnNode.GetName(), podTrafficNeighbor) {
				continue
			}

			score := s.CalculateTrafficScore(podOnNode.GetName(), currentPod)
			scores = append(scores, score)
		}
	}

	return NodeScore{
		NodeName: nodeName,
		PodName:  currentPod,
		Score:    sum(scores),
	}
}

func InitKubernetesConnection() *kubernetes.Clientset {
	// Out-of-cluster configuration
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the client set
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	return clientset
}

func (s *Scoring) ListAllNodes() *v1.NodeList {
	nodes, err := s.clientSet.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	return nodes
}

func (s *Scoring) ListAllPods() *v1.PodList {
	pods, err := s.clientSet.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	appPods := &v1.PodList{}
	for _, pod := range pods.Items {
		if isNetmarksPod(pod.Labels["app"]) {
			appPods.Items = append(appPods.Items, pod)
		}
	}

	return appPods
}

func main() {
	var scores []NodeScore

	scoringPlugin := Scoring{
		prometheus: NewPrometheusClient("http://localhost:9090"),
		clientSet:  InitKubernetesConnection(),
	}

	nodes := scoringPlugin.ListAllNodes()
	pods := scoringPlugin.ListAllPods()

	for _, pod := range pods.Items {
		for _, node := range nodes.Items {
			scores = append(scores, scoringPlugin.CalculateScore(node.GetName(), pod.GetName()))
		}
	}

	fmt.Printf("%d\n", len(scores))
	for _, s := range scores {
		fmt.Printf("Pod %s on Node %s has a score of %f\n", s.PodName, s.NodeName, s.Score)
	}
}
