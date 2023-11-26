package main

import (
	Bread "NetMARKS/services/bread/proto"
	Fish "NetMARKS/services/fish/proto"
	Gold "NetMARKS/services/gold/proto"
	Meat "NetMARKS/services/meat/proto"
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

const ServiceName = "Gold"
const ServicePort = "8080"

const BreadServiceAddr = "netmarks-bread.default.svc.cluster.local:8080"
const FishServiceAddr = "netmarks-fish.default.svc.cluster.local:8080"
const MeatServiceAddr = "netmarks-meat.default.svc.cluster.local:8080"

var NodeName = os.Getenv("K8S_NODE_NAME")
var RequestCount = shared.InitPrometheusRequestCountMetrics()

// --------------- gRPC Methods ---------------

type GoldServer struct {
	Gold.UnimplementedGoldServer
	breadClient Bread.BreadClient
	fishClient  Fish.FishClient
	meatClient  Meat.MeatClient
}

func NewGoldServer(b Bread.BreadClient, f Fish.FishClient, m Meat.MeatClient) *GoldServer {
	return &GoldServer{
		breadClient: b,
		fishClient:  f,
		meatClient:  m,
	}
}

func newGRPCServer(lis net.Listener) error {
	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	breadClient := Bread.NewBreadClient(shared.InitGrpcClientConn(BreadServiceAddr))
	fishClient := Fish.NewFishClient(shared.InitGrpcClientConn(FishServiceAddr))
	meatClient := Meat.NewMeatClient(shared.InitGrpcClientConn(MeatServiceAddr))

	Gold.RegisterGoldServer(grpcServer, NewGoldServer(breadClient, fishClient, meatClient))

	return grpcServer.Serve(lis)
}

func (s *GoldServer) Produce(ctx context.Context, req *Gold.Request) (*Gold.Response, error) {
	shared.SetGRPCHeader(&ctx)
	ctx, span := shared.InitServerSpan(ctx, ServiceName)
	defer span.End()
	defer RequestCount.With(prometheus.Labels{
		"service_name": ServiceName,
		"node_name":    NodeName,
	}).Inc()

	latency, _ := strconv.ParseInt(os.Getenv("CONSTANT_LATENCY"), 10, 32)

	var wg sync.WaitGroup
	wg.Add(3)

	ch := make(chan shared.GRPCResponse)
	var mutex sync.Mutex

	r := Gold.Response{}
	for i := uint64(0); i < req.Quantity; i++ {
		r.Quantity += 1

		product := Gold.Single{
			Id:             shared.GenerateRandomUUID(),
			RandomMetadata: shared.GenerateFakeMetadataString(ctx, req.ResponseSize),
		}

		go shared.ConcurrentGRPCBread(ctx, req.ResponseSize, s.breadClient, &wg, ch)
		go shared.ConcurrentGRPCFish(ctx, req.ResponseSize, s.fishClient, &wg, ch)
		go shared.ConcurrentGRPCMeat(ctx, req.ResponseSize, s.meatClient, &wg, ch)
		go func() {
			wg.Wait()
			close(ch)
		}()

		for response := range ch {
			mutex.Lock()
			if response.Type == "bread" {
				product.BreadId = response.Body
			} else if response.Type == "fish" {
				product.FishId = response.Body
			} else if response.Type == "meat" {
				product.MeatId = response.Body
			}
			mutex.Unlock()
		}

		r.Items = append(r.Items, &product)

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
	responseSize := r.URL.Query().Get("response_size")

	latency, _ := strconv.ParseInt(os.Getenv("CONSTANT_LATENCY"), 10, 32)

	var wg sync.WaitGroup
	wg.Add(3)

	ch := make(chan shared.HTTPResponse)
	responses := make(map[string]shared.HTTPResponse)
	var mutex sync.Mutex

	response := shared.CoalHTTPResponse{
		Type: ServiceName,
	}
	for i := uint64(0); i < quantity; i++ {
		response.Quantity += 1

		product := shared.SingleCoal{
			Id:             shared.GenerateRandomUUID(),
			RandomMetadata: shared.GenerateFakeMetadataString(ctx, responseSize),
		}

		go shared.ConcurrentHTTPRequest("http://"+BreadServiceAddr+"?response_size="+responseSize, "bread", &wg, ch)
		go shared.ConcurrentHTTPRequest("http://"+FishServiceAddr+"?response_size="+responseSize, "fish", &wg, ch)
		go shared.ConcurrentHTTPRequest("http://"+MeatServiceAddr+"?response_size="+responseSize, "meat", &wg, ch)
		go func() {
			wg.Wait()
			close(ch)
		}()

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
				if response.Type == "fish" {
					var fish shared.BasicTypeHTTPResponse
					err := json.Unmarshal(response.Body, &fish)
					if err != nil {
						w.Write([]byte(err.Error()))
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
					product.FishId = fish.Items[0].Id
				} else if response.Type == "bread" {
					var bread shared.BreadHTTPResponse
					err := json.Unmarshal(response.Body, &bread)
					if err != nil {
						w.Write([]byte(err.Error()))
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
					product.BreadId = bread.Items[0].Id
				} else if response.Type == "meat" {
					var meat shared.MeatHTTPResponse
					err := json.Unmarshal(response.Body, &meat)
					if err != nil {
						w.Write([]byte(err.Error()))
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
					product.MeatId = meat.Items[0].Id
				}
			}
		}

		response.Items = append(response.Items, product)

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
