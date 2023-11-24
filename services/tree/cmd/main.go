package main

import (
	"context"
	"encoding/json"
	"fmt"
	Tree "github.com/JBHua/NetMARKS/services/tree/proto"
	"github.com/JBHua/NetMARKS/shared"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"
)

var ServiceName = "Tree"
var ServicePortEnv = "TREE_SERVICE_PORT"

var NodeName = os.Getenv("K8S_NODE_NAME")
var RequestCount = shared.InitPrometheusRequestCountMetrics()

// --------------- gRPC Methods ---------------

type TreeServer struct {
	Tree.UnimplementedTreeServer
	logger *otelzap.SugaredLogger
}

func NewTreeServer(l *otelzap.SugaredLogger) *TreeServer {
	return &TreeServer{
		logger: l,
	}
}

func (s *TreeServer) Produce(ctx context.Context, req *Tree.Request) (*Tree.Response, error) {
	shared.SetGRPCHeader(&ctx)
	ctx, span := shared.InitServerSpan(ctx, ServiceName)
	defer span.End()

	latency, _ := strconv.ParseInt(os.Getenv("CONSTANT_LATENCY"), 10, 32)

	r := Tree.Response{}
	for i := uint64(0); i < req.Quantity; i++ {
		r.Quantity += 1
		r.Items = append(r.Items, &Tree.Single{
			Id:             shared.GenerateRandomUUID(),
			RandomMetadata: shared.GenerateFakeMetadataInByte(ctx, req.ResponseSize),
		})

		time.Sleep(time.Duration(latency) * time.Millisecond)
	}

	span.SetStatus(codes.Ok, "success")
	return &r, nil
}

// --------------- HTTP Methods ---------------

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

	var responseSize uint64
	responseSize, err = strconv.ParseUint(r.URL.Query().Get("response_size"), 10, 64)
	if err != nil {
		responseSize = 1
	}

	latency, _ := strconv.ParseInt(os.Getenv("CONSTANT_LATENCY"), 10, 32)

	response := shared.BasicTypeHTTPResponse{
		Type: ServiceName,
	}
	for i := uint64(0); i < quantity; i++ {
		response.Quantity += 1
		response.Items = append(response.Items, shared.SingleBasicType{
			Id:             shared.GenerateRandomUUID(),
			RandomMetadata: shared.GenerateFakeMetadataInByte(ctx, responseSize),
		})

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

	useGRPC, _ := strconv.ParseBool(os.Getenv("USE_GRPC"))
	if useGRPC {
		logger.Info("Using GRPC")
		listener, err := net.Listen("tcp", fmt.Sprintf(":%s", os.Getenv(ServicePortEnv)))
		if err != nil {
			logger.Fatalf("could not attach listener to port: %v. %v", os.Getenv(ServicePortEnv), err)
		}
		logger.Infof("Running at %s\n", os.Getenv(ServicePortEnv))

		grpcServer := grpc.NewServer()
		Tree.RegisterTreeServer(grpcServer, NewTreeServer(logger))

		go func() {
			if err := grpcServer.Serve(listener); err != nil {
				logger.Fatalf("could not start grpc server: %v", err)
			}
		}()

		shared.MonitorShutdownSignal()
	} else {
		logger.Info("Using HTTP")
		mux := http.NewServeMux()
		mux.HandleFunc("/", Produce)

		// Start HTTP Server
		port := os.Getenv(ServicePortEnv)
		logger.Info(port)
		err := http.ListenAndServe(":"+port, mux)

		if err != nil {
			panic(err)
		}
		logger.Infof("service running on port %s\n", port)
	}
}
