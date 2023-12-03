package main

import (
	Coal "NetMARKS/services/coal/proto"
	Coin "NetMARKS/services/coin/proto"
	Gold "NetMARKS/services/gold/proto"
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

const ServiceName = "Coin"
const ServicePort = "8080"

const CoalServiceAddr = "netmarks-coal.default.svc.cluster.local:8080"
const GoldServiceAddr = "netmarks-gold.default.svc.cluster.local:8080"

var NodeName = os.Getenv("K8S_NODE_NAME")
var RequestCount, InterNodeRequestCount = shared.InitPrometheusRequestCountMetrics()

// --------------- gRPC Methods ---------------

type CoinServer struct {
	Coin.UnimplementedCoinServer
	coalClient Coal.CoalClient
	goldClient Gold.GoldClient
}

func NewCoinServer(c Coal.CoalClient, g Gold.GoldClient) *CoinServer {
	return &CoinServer{
		coalClient: c,
		goldClient: g,
	}
}

func newGRPCServer(lis net.Listener) error {
	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	coalClient := Coal.NewCoalClient(shared.InitGrpcClientConn(GoldServiceAddr))
	goldClient := Gold.NewGoldClient(shared.InitGrpcClientConn(CoalServiceAddr))
	Coin.RegisterCoinServer(grpcServer, NewCoinServer(coalClient, goldClient))

	return grpcServer.Serve(lis)
}

func (s *CoinServer) Produce(ctx context.Context, req *Coin.Request) (*Coin.Response, error) {
	shared.SetGRPCHeader(&ctx)
	ctx, span := shared.InitServerSpan(ctx, ServiceName)
	defer span.End()

	var wg sync.WaitGroup
	wg.Add(2)

	ch := make(chan shared.GRPCResponse)
	var mutex sync.Mutex

	r := Coin.Response{}
	for i := uint64(0); i < req.Quantity; i++ {
		r.Quantity += 1

		singleCoin := Coin.Single{
			Id:             shared.GenerateRandomUUID(),
			RandomMetadata: shared.GenerateFakeMetadataString(ctx, req.ResponseSize),
		}

		go shared.ConcurrentGRPCCoal(ctx, req.ResponseSize, s.coalClient, &wg, ch)
		go shared.ConcurrentGRPCGold(ctx, req.ResponseSize, s.goldClient, &wg, ch)
		go func() {
			wg.Wait()
			close(ch)
		}()

		for response := range ch {
			mutex.Lock()
			if response.Type == "coal" {
				singleCoin.CoalId = response.Body
			} else if response.Type == "gold" {
				singleCoin.GoldId = response.Body
			}
			mutex.Unlock()
		}

		r.Items = append(r.Items, &singleCoin)
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

	r.WithContext(ctx)
	w.Header().Set("Content-Type", "application/json")

	var quantity uint64
	quantity, err := strconv.ParseUint(r.URL.Query().Get("quantity"), 10, 64)
	if err != nil {
		quantity = 1
	}
	responseSize := r.URL.Query().Get("response_size")

	var wg sync.WaitGroup
	wg.Add(2)

	ch := make(chan shared.HTTPResponse)
	responses := make(map[string]shared.HTTPResponse)
	var mutex sync.Mutex

	response := shared.CoinHTTPResponse{
		Type: ServiceName,
	}
	for i := uint64(0); i < quantity; i++ {
		response.Quantity += 1

		latency := shared.CalculateArtificialLatency(r.Header, NodeName)
		requestId, originalRequestService, upstreamNodeName := shared.ExtractUpstreamRequestID(r.Header, ServiceName, NodeName)
		singleCoin := shared.SingleCoin{
			Id:             requestId,
			RandomMetadata: shared.GenerateFakeMetadataString(ctx, responseSize),
		}

		go shared.ConcurrentHTTPRequest("http://"+CoalServiceAddr+"?response_size="+responseSize, "coal", NodeName, requestId, originalRequestService, &wg, ch)
		go shared.ConcurrentHTTPRequest("http://"+GoldServiceAddr+"?response_size="+responseSize, "gold", NodeName, requestId, originalRequestService, &wg, ch)
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
				if response.Type == "coal" {
					var coal shared.CoalHTTPResponse
					err := json.Unmarshal(response.Body, &coal)
					if err != nil {
						w.Write([]byte(err.Error()))
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
					singleCoin.CoalId = coal.Items[0].Id
				} else if response.Type == "gold" {
					var gold shared.CoalHTTPResponse
					err := json.Unmarshal(response.Body, &gold)
					if err != nil {
						w.Write([]byte(err.Error()))
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
					singleCoin.GoldId = gold.Items[0].Id
				}
			}
		}

		response.Items = append(response.Items, singleCoin)

		shared.UpdateRequestMetrics(RequestCount, InterNodeRequestCount, originalRequestService, ServiceName, NodeName, upstreamNodeName)
		time.Sleep(time.Duration(latency) * time.Millisecond)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// --------------- Main Logic ---------------

func main() {
	logger := shared.InitSugaredLogger()
	shared.ConfigureRuntime()
	http.DefaultTransport.(*http.Transport).MaxIdleConnsPerHost = 200
	prometheus.MustRegister(RequestCount)
	prometheus.MustRegister(InterNodeRequestCount)

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
