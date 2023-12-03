package main

import (
	"context"
	"fmt"
	"github.com/prometheus/client_golang/api"
	promethus_v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"net/http"
	"sort"
	"strconv"
	"sync"
	"time"
)

const responseSize = "4k"
const iterationCount = 50

var services = map[string]string{
	"Beer":    "63583",
	"Board":   "63586",
	"Boat":    "50773",
	"Bread":   "63590",
	"Coal":    "63592",
	"Coin":    "63594",
	"Flour":   "63598",
	"Gold":    "63600",
	"Iron":    "63606",
	"Ironore": "63608",
	"Log":     "63610",
	"Meat":    "63612",
	"Pig":     "63614",
	"Sword":   "63616",
	"Tools":   "63618",
}

const (
	nodeMesaureQuery = "sum_over_time(process_cpu_seconds_total[1y])"
	// IstioResponseByteSum TODO: add time range parameter
	SubRequestCount          = "sub_request_count_per_service{original_request_service=\"%s\"}"
	InterNodeSubRequestCount = "inter_node_sub_request_count_per_service{original_request_service=\"%s\"}"
)

func GetSubRequestCountByService(serviceName string, internode bool, client promethus_v1.API) int64 {
	var total int64 = 0

	var query string
	if internode {
		query = InterNodeSubRequestCount
	} else {
		query = SubRequestCount
	}

	value, _, err := client.Query(context.Background(), fmt.Sprintf(query, serviceName), time.Now())
	if err != nil {
		panic(err)
	}

	nodeMeasure := value.(model.Vector)
	for _, res := range nodeMeasure {
		if s, err := strconv.ParseInt(res.Value.String(), 10, 64); err == nil {
			total += s
		}
	}

	return total
}

func MakeHTTPRequest(serviceName, port string, wg *sync.WaitGroup) int {
	defer wg.Done()
	println("Sending request to: " + serviceName)
	url := "http://127.0.0.1:" + port + "?response_size=" + responseSize
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	return resp.StatusCode
}

func main() {
	// Step 1: Iterate over all services, sending each service 10 requests
	var wg sync.WaitGroup
	for serviceName, servicePort := range services {
		for i := 0; i < iterationCount; i++ {
			wg.Add(1)
			go MakeHTTPRequest(serviceName, servicePort, &wg)
			time.Sleep(100 * time.Millisecond)
		}

		wg.Wait()
	}
	// Wait till metrics got scrapped by prometheus
	time.Sleep(30 * time.Second)

	// Step 2: For every service, add up the total sub-request count && internode sub-request count
	client, err := api.NewClient(api.Config{Address: "http://localhost:9090"})
	if err != nil {
		panic(err)
	}

	result := make(map[string]int64)
	promoClient := promethus_v1.NewAPI(client)
	for serviceName, _ := range services {
		subRequestCount := GetSubRequestCountByService(serviceName, false, promoClient)
		result[serviceName] += subRequestCount

		InterNodeSubRequestCount := GetSubRequestCountByService(serviceName, true, promoClient)
		result["InterNode"+serviceName] += InterNodeSubRequestCount
	}

	keys := make([]string, 0)
	for k, _ := range services {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, serviceName := range keys {
		fmt.Printf("Serivice: %s has a total sub-request count of: %f\n", serviceName, float64(result[serviceName])/iterationCount)
		fmt.Printf("Serivice: %s has a total InterNode request of: %f\n", serviceName, float64(result["InterNode"+serviceName])/iterationCount)

		fmt.Printf("\n")
	}
}
