package main

import (
	Board "NetMARKS/services/board/proto"
	Iron "NetMARKS/services/iron/proto"
	Tools "NetMARKS/services/tools/proto"
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

const ServiceName = "Tools"
const ServicePort = "8080"

const BoardServiceAddr = "netmarks-board.default.svc.cluster.local:8080"
const IronServiceAddr = "netmarks-iron.default.svc.cluster.local:8080"

var NodeName = os.Getenv("K8S_NODE_NAME")
var RequestCount, InterNodeRequestCount = shared.InitPrometheusRequestCountMetrics()

// --------------- gRPC Methods ---------------

type ToolsServer struct {
	Tools.UnimplementedToolsServer
	boardClient Board.BoardClient
	ironClient  Iron.IronClient
}

func NewToolsServer(b Board.BoardClient, i Iron.IronClient) *ToolsServer {
	return &ToolsServer{
		boardClient: b,
		ironClient:  i,
	}
}

func newGRPCServer(lis net.Listener) error {
	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	boardClient := Board.NewBoardClient(shared.InitGrpcClientConn(BoardServiceAddr))
	ironClient := Iron.NewIronClient(shared.InitGrpcClientConn(IronServiceAddr))
	Tools.RegisterToolsServer(grpcServer, NewToolsServer(boardClient, ironClient))

	return grpcServer.Serve(lis)
}

func (s *ToolsServer) Produce(ctx context.Context, req *Tools.Request) (*Tools.Response, error) {
	shared.SetGRPCHeader(&ctx)
	ctx, span := shared.InitServerSpan(ctx, ServiceName)
	defer span.End()

	var wg sync.WaitGroup
	wg.Add(2)

	ch := make(chan shared.GRPCResponse)
	var mutex sync.Mutex

	r := Tools.Response{}
	for i := uint64(0); i < req.Quantity; i++ {
		r.Quantity += 1

		singleTool := Tools.Single{
			Id:             shared.GenerateRandomUUID(),
			RandomMetadata: shared.GenerateFakeMetadataString(ctx, req.ResponseSize),
		}

		go shared.ConcurrentGRPCBoard(ctx, req.ResponseSize, s.boardClient, &wg, ch)
		go shared.ConcurrentGRPCIron(ctx, req.ResponseSize, s.ironClient, &wg, ch)
		go func() {
			wg.Wait()
			close(ch)
		}()

		for response := range ch {
			mutex.Lock()
			if response.Type == "board" {
				singleTool.BoardId = response.Body
			} else if response.Type == "ironore" {
				singleTool.IronoreId = response.Body
			}
			mutex.Unlock()
		}

		r.Items = append(r.Items, &singleTool)
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

	response := shared.ToolsHTTPResponse{
		Type: ServiceName,
	}
	for i := uint64(0); i < quantity; i++ {
		response.Quantity += 1

		latency := shared.CalculateArtificialLatency(r.Header, NodeName)
		requestId, originalRequestService, upstreamNodeName := shared.ExtractUpstreamRequestID(r.Header, ServiceName, NodeName)
		singleTool := shared.SingleTool{
			Id:             requestId,
			RandomMetadata: shared.GenerateFakeMetadataString(ctx, responseSize),
		}

		go shared.ConcurrentHTTPRequest("http://"+BoardServiceAddr+"?response_size="+responseSize, "board", NodeName, requestId, originalRequestService, &wg, ch)
		go shared.ConcurrentHTTPRequest("http://"+IronServiceAddr+"?response_size="+responseSize, "ironore", NodeName, requestId, originalRequestService, &wg, ch)
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
				if response.Type == "board" {
					var board shared.BoardHTTPResponse
					err := json.Unmarshal(response.Body, &board)
					if err != nil {
						w.Write([]byte(err.Error()))
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
					singleTool.BoardId = board.Items[0].Id
				} else if response.Type == "ironore" {
					var ironore shared.CoalHTTPResponse
					err := json.Unmarshal(response.Body, &ironore)
					if err != nil {
						w.Write([]byte(err.Error()))
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
					singleTool.IronoreId = ironore.Items[0].Id
				}
			}
		}

		response.Items = append(response.Items, singleTool)

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
