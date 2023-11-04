package main

import (
	"NetMARKS/shared"
	"context"
	"encoding/json"
	"fmt"
	Water "github.com/JBHua/NetMARKS/services/water/proto"
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

var ServiceName = "Water"
var ServicePortEnv = "Water_SERVICE_PORT"

// --------------- gRPC Methods ---------------

type WaterServer struct {
	Water.UnimplementedWaterServer
	logger *otelzap.SugaredLogger
}

func NewWaterServer(l *otelzap.SugaredLogger) *WaterServer {
	return &WaterServer{
		logger: l,
	}
}

func (s *WaterServer) ProduceWater(ctx context.Context, req *Water.Request) (*Water.Single, error) {
	shared.SetGRPCHeader(&ctx)
	ctx, span := shared.InitServerSpan(ctx, ServiceName)
	defer span.End()

	latency, _ := strconv.ParseInt(os.Getenv("CONSTANT_LATENCY"), 10, 32)
	time.Sleep(time.Duration(latency) * time.Millisecond)

	span.SetStatus(codes.Ok, "success")

	t := &Water.Single{
		Id:             shared.GenerateRandomUUID(),
		RandomMetadata: shared.GenerateFakeMetadata(),
	}

	return t, nil
}

// --------------- HTTP Methods ---------------

func ProduceWater(w http.ResponseWriter, r *http.Request) {
	ctx, span := shared.InitServerSpan(context.Background(), ServiceName)
	defer span.End()

	r.WithContext(ctx)
	w.Header().Set("Content-Type", "application/json")

	latency, _ := strconv.ParseInt(os.Getenv("CONSTANT_LATENCY"), 10, 32)
	time.Sleep(time.Duration(latency) * time.Millisecond)

	d := shared.WaterHTTP{
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
		Water.RegisterWaterServer(grpcServer, NewWaterServer(logger))

		go func() {
			if err := grpcServer.Serve(listener); err != nil {
				logger.Fatalf("could not start grpc server: %v", err)
			}
		}()

		shared.MonitorShutdownSignal()
	} else {
		logger.Info("Using HTTP")
		mux := http.NewServeMux()
		mux.HandleFunc("/", ProduceWater)

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
