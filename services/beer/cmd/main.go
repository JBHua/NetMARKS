package main

import (
	Beer "NetMARKS/services/beer/proto"
	Grain "NetMARKS/services/grain/proto"
	Water "NetMARKS/services/water/proto"
	"NetMARKS/shared"
	"context"
	"encoding/json"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/soheilhy/cmux"
	"go.opentelemetry.io/otel/codes"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

const ServiceName = "Beer"
const ServicePort = "8080"

// const GrainServiceAddr = "netmarks-grain.default.svc.cluster.local:8080"
// const WaterServiceAddr = "netmarks-water.default.svc.cluster.local:8080"
const GrainServiceAddr = "127.0.0.1:55292"
const WaterServiceAddr = "127.0.0.1:55298"

var NodeName = os.Getenv("K8S_NODE_NAME")
var RequestCount = shared.InitPrometheusRequestCountMetrics()

// --------------- gRPC Methods ---------------

type BeerServer struct {
	Beer.UnimplementedBeerServer
	waterClient Water.WaterClient
	grainClient Grain.GrainClient
}

func NewBeerServer(w Water.WaterClient, g Grain.GrainClient) *BeerServer {
	return &BeerServer{
		waterClient: w,
		grainClient: g,
	}
}

func newGRPCServer(lis net.Listener) error {
	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	waterClient := Water.NewWaterClient(shared.InitGrpcClientConn(WaterServiceAddr))
	grainClient := Grain.NewGrainClient(shared.InitGrpcClientConn(GrainServiceAddr))
	Beer.RegisterBeerServer(grpcServer, NewBeerServer(waterClient, grainClient))

	return grpcServer.Serve(lis)
}

func (s *BeerServer) Produce(ctx context.Context, req *Beer.Request) (*Beer.Response, error) {
	shared.SetGRPCHeader(&ctx)
	ctx, span := shared.InitServerSpan(ctx, ServiceName)
	defer span.End()
	defer RequestCount.With(prometheus.Labels{
		"service_name": ServiceName,
		"node_name":    NodeName,
	}).Inc()

	latency, _ := strconv.ParseInt(os.Getenv("CONSTANT_LATENCY"), 10, 32)

	var wg sync.WaitGroup
	wg.Add(2)

	ch := make(chan shared.ConcurrentGRPCResponse)
	var mutex sync.Mutex

	r := Beer.Response{}
	for i := uint64(0); i < req.Quantity; i++ {
		r.Quantity += 1

		singleBeer := Beer.Single{
			Id:             shared.GenerateRandomUUID(),
			RandomMetadata: shared.GenerateFakeMetadataString(ctx, req.ResponseSize),
		}

		go shared.ConcurrentGRPCWater(ctx, s.waterClient, &wg, ch)
		go shared.ConcurrentGRPCGrain(ctx, s.grainClient, &wg, ch)
		go func() {
			wg.Wait()
			close(ch)
		}()

		for response := range ch {
			mutex.Lock()
			if response.Type == "water" {
				singleBeer.WaterId = response.Body
			} else if response.Type == "grain" {
				singleBeer.GrainId = response.Body
			}
			mutex.Unlock()
		}

		r.Items = append(r.Items, &singleBeer)

		time.Sleep(time.Duration(latency) * time.Millisecond)
	}

	span.SetStatus(codes.Ok, "success")
	return &r, nil
}

// --------------- HTTP Methods ---------------

func newHTTPServer(lis net.Listener) error {
	mux := http.NewServeMux()

	mux.HandleFunc("/", Produce)
	mux.Handle("/metrics", promhttp.Handler())

	s := &http.Server{Handler: mux}
	return s.Serve(lis)
}

func Produce(w http.ResponseWriter, r *http.Request) {
	ctx, span := shared.InitServerSpan(context.Background(), ServiceName)
	defer span.End()
	defer RequestCount.With(prometheus.Labels{
		"service_name": ServiceName,
		"node_name":    NodeName,
	}).Inc()

	r.WithContext(ctx)
	w.Header().Set("Content-Type", "application/json")

	var quantity uint64
	quantity, err := strconv.ParseUint(r.URL.Query().Get("quantity"), 10, 64)
	if err != nil {
		quantity = 1
	}

	latency, _ := strconv.ParseInt(os.Getenv("CONSTANT_LATENCY"), 10, 32)

	var wg sync.WaitGroup
	wg.Add(2)

	ch := make(chan shared.HTTPResponse)
	responses := make(map[string]shared.HTTPResponse)
	var mutex sync.Mutex

	response := shared.BeerHTTPResponse{
		Type: ServiceName,
	}
	for i := uint64(0); i < quantity; i++ {
		response.Quantity += 1

		singleBeer := shared.SingleBeer{
			Id:             shared.GenerateRandomUUID(),
			RandomMetadata: shared.GenerateFakeMetadataString(ctx, r.URL.Query().Get("response_size")),
		}

		go shared.ConcurrentHTTPRequest("http://"+GrainServiceAddr, "grain", &wg, ch)
		go shared.ConcurrentHTTPRequest("http://"+WaterServiceAddr, "water", &wg, ch)
		go func() {
			wg.Wait()
			close(ch)
		}()

		//for response := range ch {
		//	mutex.Lock()
		//	responses[response.Type] = response
		//	mutex.Unlock()
		//}

		for response := range ch {
			mutex.Lock()
			responses[response.Type] = response
			mutex.Unlock()
		}

		for url, response := range responses {
			if response.Err != nil {
				fmt.Printf("Error fetching %s: %v\n", url, response.Err)
				w.Write([]byte(err.Error()))
				w.WriteHeader(http.StatusInternalServerError)
				return
			} else {
				fmt.Printf("Response from %s: %s\n", url, response.Body)
				if response.Type == "water" {
					var water shared.BasicTypeHTTPResponse
					err := json.Unmarshal(response.Body, &water)
					if err != nil {
						w.Write([]byte(err.Error()))
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
					singleBeer.WaterId = water.Items[0].Id
				} else if response.Type == "grain" {
					var grain shared.BasicTypeHTTPResponse
					err := json.Unmarshal(response.Body, &grain)
					if err != nil {
						w.Write([]byte(err.Error()))
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
					singleBeer.GrainId = grain.Items[0].Id
				}
			}
		}

		response.Items = append(response.Items, singleBeer)

		time.Sleep(time.Duration(latency) * time.Millisecond)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// --------------- Main Logic ---------------

func main() {
	logger := shared.InitSugaredLogger()
	shared.ConfigureRuntime()
	prometheus.MustRegister(RequestCount)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", ServicePort))
	if err != nil {
		logger.Fatalf("could not attach listener to port: %v. %v", ServicePort, err)
	}

	mux := cmux.New(listener)
	grpcListener := mux.Match(cmux.HTTP2HeaderField("content-type", "application/grpc"))
	httpListener := mux.Match(cmux.HTTP1Fast())

	// Use an error group to start all of them
	g := errgroup.Group{}
	g.Go(func() error { return newGRPCServer(grpcListener) })
	g.Go(func() error { return newHTTPServer(httpListener) })
	g.Go(func() error { return mux.Serve() })

	log.Println("run server:", g.Wait())
}