package main

import (
	Coal "NetMARKS/services/coal/proto"
	Iron "NetMARKS/services/iron/proto"
	Sword "NetMARKS/services/sword/proto"
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

const ServiceName = "Sword"
const ServicePort = "8080"

const CoalServiceAddr = "netmarks-coal.default.svc.cluster.local:8080"
const IronServiceAddr = "netmarks-iron.default.svc.cluster.local:8080"

var NodeName = os.Getenv("K8S_NODE_NAME")
var RequestCount, InterNodeRequestCount = shared.InitPrometheusRequestCountMetrics()

// --------------- gRPC Methods ---------------

type SwordServer struct {
	Sword.UnimplementedSwordServer
	coalClient Coal.CoalClient
	ironClient Iron.IronClient
}

func NewSwordServer(c Coal.CoalClient, i Iron.IronClient) *SwordServer {
	return &SwordServer{
		coalClient: c,
		ironClient: i,
	}
}

func newGRPCServer(lis net.Listener) error {
	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	coalClient := Coal.NewCoalClient(shared.InitGrpcClientConn(CoalServiceAddr))
	ironClient := Iron.NewIronClient(shared.InitGrpcClientConn(IronServiceAddr))
	Sword.RegisterSwordServer(grpcServer, NewSwordServer(coalClient, ironClient))

	return grpcServer.Serve(lis)
}

func (s *SwordServer) Produce(ctx context.Context, req *Sword.Request) (*Sword.Response, error) {
	shared.SetGRPCHeader(&ctx)
	ctx, span := shared.InitServerSpan(ctx, ServiceName)
	defer span.End()

	latency, _ := strconv.ParseInt(os.Getenv("CONSTANT_LATENCY"), 10, 32)

	var wg sync.WaitGroup
	wg.Add(2)

	ch := make(chan shared.GRPCResponse)
	var mutex sync.Mutex

	r := Sword.Response{}
	for i := uint64(0); i < req.Quantity; i++ {
		r.Quantity += 1

		singleSword := Sword.Single{
			Id:             shared.GenerateRandomUUID(),
			RandomMetadata: shared.GenerateFakeMetadataString(ctx, req.ResponseSize),
		}

		go shared.ConcurrentGRPCCoal(ctx, req.ResponseSize, s.coalClient, &wg, ch)
		go shared.ConcurrentGRPCIron(ctx, req.ResponseSize, s.ironClient, &wg, ch)
		go func() {
			wg.Wait()
			close(ch)
		}()

		for response := range ch {
			mutex.Lock()
			if response.Type == "coal" {
				singleSword.CoalId = response.Body
			} else if response.Type == "iron" {
				singleSword.IronId = response.Body
			}
			mutex.Unlock()
		}

		r.Items = append(r.Items, &singleSword)

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
	wg.Add(2)

	ch := make(chan shared.HTTPResponse)
	responses := make(map[string]shared.HTTPResponse)
	var mutex sync.Mutex

	response := shared.SwordHTTPResponse{
		Type: ServiceName,
	}
	for i := uint64(0); i < quantity; i++ {
		response.Quantity += 1

		requestId, originalRequestService, upstreamNodeName := shared.ExtractUpstreamRequestID(r.Header, ServiceName, NodeName)
		singleSword := shared.SingleSword{
			Id:             requestId,
			RandomMetadata: shared.GenerateFakeMetadataString(ctx, responseSize),
		}

		go shared.ConcurrentHTTPRequest("http://"+CoalServiceAddr+"?response_size="+responseSize, "coal", NodeName, requestId, originalRequestService, &wg, ch)
		go shared.ConcurrentHTTPRequest("http://"+IronServiceAddr+"?response_size="+responseSize, "iron", NodeName, requestId, originalRequestService, &wg, ch)
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
					singleSword.CoalId = coal.Items[0].Id
				} else if response.Type == "iron" {
					var iron shared.CoalHTTPResponse
					err := json.Unmarshal(response.Body, &iron)
					if err != nil {
						w.Write([]byte(err.Error()))
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
					singleSword.IronId = iron.Items[0].Id
				}
			}
		}

		response.Items = append(response.Items, singleSword)

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
