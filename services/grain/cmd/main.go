package main

import (
	"NetMARKS/shared"
	"context"
	"encoding/json"
	"fmt"
	Grain "github.com/JBHua/NetMARKS/services/grain/proto"
	"github.com/joho/godotenv"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"
)

var ServiceName = "Grain"
var ServicePortEnv = "GRAIN_SERVICE_PORT"

// --------------- gRPC Methods ---------------

type GrainServer struct {
	Grain.UnimplementedGrainServer
	logger *otelzap.SugaredLogger
}

func NewGrainServer(l *otelzap.SugaredLogger) *GrainServer {
	return &GrainServer{
		logger: l,
	}
}

func (s *GrainServer) ProduceGrain(ctx context.Context, req *Grain.Request) (*Grain.Single, error) {
	shared.SetGRPCHeader(&ctx)
	ctx, span := shared.InitServerSpan(ctx, ServiceName)
	defer span.End()

	latency, _ := strconv.ParseInt(os.Getenv("CONSTANT_LATENCY"), 10, 32)
	time.Sleep(time.Duration(latency) * time.Millisecond)

	span.SetStatus(codes.Ok, "success")

	g := &Grain.Single{
		Id:             shared.GenerateRandomUUID(),
		RandomMetadata: shared.GenerateFakeMetadata(),
	}

	return g, nil
}

// --------------- HTTP Methods ---------------

func ProduceGrain(w http.ResponseWriter, r *http.Request) {
	ctx, span := shared.InitServerSpan(context.Background(), ServiceName)
	defer span.End()

	r.WithContext(ctx)
	w.Header().Set("Content-Type", "application/json")

	latency, _ := strconv.ParseInt(os.Getenv("CONSTANT_LATENCY"), 10, 32)
	time.Sleep(time.Duration(latency) * time.Millisecond)

	d := shared.GrainHTTP{
		Id:             shared.GenerateRandomUUID(),
		RandomMetadata: shared.GenerateFakeMetadata(),
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(d)
}

// --------------- Main Logic ---------------

func main() {
	logger := shared.InitSugaredLogger()

	err := godotenv.Load("../../shared/.env")
	if err != nil {
		logger.Fatal("Error loading .env file")
	}

	shared.ConfigureRuntime()

	useGRPC, _ := strconv.ParseBool(os.Getenv("USE_GRPC"))
	if useGRPC {
		logger.Info("Using GRPC")
		listener, err := net.Listen("tcp", fmt.Sprintf(":%s", os.Getenv(ServicePortEnv)))
		if err != nil {
			logger.Fatalf("could not attach listener to port: %v. %v", os.Getenv(ServicePortEnv), err)
		}
		logger.Infof("Running at %s\n", os.Getenv(ServicePortEnv))

		grpcServer := grpc.NewServer()
		Grain.RegisterGrainServer(grpcServer, NewGrainServer(logger))

		go func() {
			if err := grpcServer.Serve(listener); err != nil {
				logger.Fatalf("could not start grpc server: %v", err)
			}
		}()

		shared.MonitorShutdownSignal()
	} else {
		logger.Info("Using HTTP")
		mux := http.NewServeMux()
		mux.HandleFunc("/", ProduceGrain)

		// Start HTTP Server
		port := os.Getenv(ServicePortEnv)
		logger.Info(port)
		err = http.ListenAndServe(":"+port, mux)

		if err != nil {
			panic(err)
		}
		logger.Infof("service running on port %s\n", port)
	}
}
