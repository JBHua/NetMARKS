package main

import (
	Coal "NetMARKS/services/coal/proto"
	Iron "NetMARKS/services/iron/proto"
	Ironore "NetMARKS/services/ironore/proto"
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

const ServiceName = "Iron"
const ServicePort = "8080"

const CoalServiceAddr = "netmarks-coal.default.svc.cluster.local:8080"
const IronoreServiceAddr = "netmarks-ironore.default.svc.cluster.local:8080"

var NodeName = os.Getenv("K8S_NODE_NAME")
var RequestCount = shared.InitPrometheusRequestCountMetrics()

// --------------- gRPC Methods ---------------

type IronServer struct {
	Iron.UnimplementedIronServer
	coalClient    Coal.CoalClient
	ironoreClient Ironore.IronoreClient
}

func NewIronServer(c Coal.CoalClient, i Ironore.IronoreClient) *IronServer {
	return &IronServer{
		coalClient:    c,
		ironoreClient: i,
	}
}

func newGRPCServer(lis net.Listener) error {
	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	coalClient := Coal.NewCoalClient(shared.InitGrpcClientConn(CoalServiceAddr))
	ironoreClient := Ironore.NewIronoreClient(shared.InitGrpcClientConn(IronoreServiceAddr))
	Iron.RegisterIronServer(grpcServer, NewIronServer(coalClient, ironoreClient))

	return grpcServer.Serve(lis)
}

func (s *IronServer) Produce(ctx context.Context, req *Iron.Request) (*Iron.Response, error) {
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

	ch := make(chan shared.GRPCResponse)
	var mutex sync.Mutex

	r := Iron.Response{}
	for i := uint64(0); i < req.Quantity; i++ {
		r.Quantity += 1

		singleIron := Iron.Single{
			Id:             shared.GenerateRandomUUID(),
			RandomMetadata: shared.GenerateFakeMetadataString(ctx, req.ResponseSize),
		}

		go shared.ConcurrentGRPCCoal(ctx, req.ResponseSize, s.coalClient, &wg, ch)
		go shared.ConcurrentGRPCIronore(ctx, req.ResponseSize, s.ironoreClient, &wg, ch)
		go func() {
			wg.Wait()
			close(ch)
		}()

		for response := range ch {
			mutex.Lock()
			if response.Type == "coal" {
				singleIron.CoalId = response.Body
			} else if response.Type == "ironore" {
				singleIron.IronoreId = response.Body
			}
			mutex.Unlock()
		}

		r.Items = append(r.Items, &singleIron)

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
	wg.Add(2)

	ch := make(chan shared.HTTPResponse)
	responses := make(map[string]shared.HTTPResponse)
	var mutex sync.Mutex

	response := shared.IronHTTPResponse{
		Type: ServiceName,
	}
	for i := uint64(0); i < quantity; i++ {
		response.Quantity += 1

		requestId, originalRequestService, upstreamNodeName := shared.ExtractUpstreamRequestID(r.Header, ServiceName, NodeName)
		singleIron := shared.SingleIron{
			Id:             requestId,
			RandomMetadata: shared.GenerateFakeMetadataString(ctx, responseSize),
		}

		go shared.ConcurrentHTTPRequest("http://"+CoalServiceAddr+"?response_size="+responseSize, "coal", NodeName, requestId, originalRequestService, &wg, ch)
		go shared.ConcurrentHTTPRequest("http://"+IronoreServiceAddr+"?response_size="+responseSize, "ironore", NodeName, requestId, originalRequestService, &wg, ch)
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
					singleIron.CoalId = coal.Items[0].Id
				} else if response.Type == "ironore" {
					var ironore shared.CoalHTTPResponse
					err := json.Unmarshal(response.Body, &ironore)
					if err != nil {
						w.Write([]byte(err.Error()))
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
					singleIron.IronoreId = ironore.Items[0].Id
				}
			}
		}

		response.Items = append(response.Items, singleIron)

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
