package main

import (
	"context"
	"encoding/json"
	"fmt"
	Flour "github.com/JBHua/NetMARKS/services/flour/proto"
	Grain "github.com/JBHua/NetMARKS/services/grain/proto"
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

const ServiceName = "Flour"
const ServicePort = "8080"
const GrainServiceAddr = "netmarks-grain.default.svc.cluster.local:5018"

var NodeName = os.Getenv("K8S_NODE_NAME")
var RequestCount = shared.InitPrometheusRequestCountMetrics()

// --------------- gRPC Methods ---------------

type FlourServer struct {
	Flour.UnimplementedFlourServer
	logger      *otelzap.SugaredLogger
	grainClient Grain.GrainClient
}

func NewFlourServer(l *otelzap.SugaredLogger, c Grain.GrainClient) *FlourServer {
	return &FlourServer{
		logger:      l,
		grainClient: c,
	}
}

func (s *FlourServer) Produce(ctx context.Context, req *Flour.Request) (*Flour.Response, error) {
	shared.SetGRPCHeader(&ctx)
	ctx, span := shared.InitServerSpan(ctx, ServiceName)
	defer span.End()

	latency, _ := strconv.ParseInt(os.Getenv("CONSTANT_LATENCY"), 10, 32)

	r := Flour.Response{}
	for i := uint64(0); i < req.Quantity; i++ {
		r.Quantity += 1

		singleGrain, err := s.grainClient.Produce(ctx, &Grain.Request{
			Quantity:     1,
			ResponseSize: 1,
		})
		if err != nil {
			return nil, err
		}

		r.Items = append(r.Items, &Flour.Single{
			Id:             shared.GenerateRandomUUID(),
			RandomMetadata: shared.GenerateFakeMetadataInByte(ctx, req.ResponseSize),
			GrainId:        singleGrain.Items[0].Id,
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

	response := shared.FlourHTTPResponse{
		Type: ServiceName,
	}
	for i := uint64(0); i < quantity; i++ {
		response.Quantity += 1

		getRes, err := http.Get("http://" + GrainServiceAddr + "?quantity=1&response_size=1")
		if err != nil {
			w.Write([]byte(err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer getRes.Body.Close()

		var grain shared.BasicTypeHTTPResponse
		err = json.NewDecoder(getRes.Body).Decode(&grain)
		if err != nil {
			w.Write([]byte(err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		response.Items = append(response.Items, shared.SingleFour{
			Id:             shared.GenerateRandomUUID(),
			RandomMetadata: shared.GenerateFakeMetadataInByte(ctx, responseSize),
			GrainId:        grain.Items[0].Id,
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
		listener, err := net.Listen("tcp", fmt.Sprintf(":%s", ServicePort))
		if err != nil {
			logger.Fatalf("could not attach listener to port: %v. %v", ServicePort, err)
		}
		logger.Infof("Running at %s\n", ServicePort)

		grainClient := Grain.NewGrainClient(shared.InitGrpcClientConn(GrainServiceAddr))

		grpcServer := grpc.NewServer()
		Flour.RegisterFlourServer(grpcServer, NewFlourServer(logger, grainClient))

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
		logger.Info("Running at %s\n", ServicePort)
		err := http.ListenAndServe(":"+ServicePort, mux)

		if err != nil {
			panic(err)
		}
		logger.Infof("service running on port %s\n", ServicePort)
	}
}
