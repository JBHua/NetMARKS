package main

import (
	"NetMARKS/shared"
	"context"
	"encoding/json"
	"fmt"
	Fish "github.com/JBHua/NetMARKS/services/fish/proto"
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

var ServiceName = "Fish"
var ServicePortEnv = "FISH_SERVICE_PORT"

// --------------- gRPC Methods ---------------

type FishServer struct {
	Fish.UnimplementedFishServer
	logger *otelzap.SugaredLogger
}

func NewFishServer(l *otelzap.SugaredLogger) *FishServer {
	return &FishServer{
		logger: l,
	}
}

func (s *FishServer) ProduceFish(ctx context.Context, req *Fish.Request) (*Fish.Single, error) {
	ctx, span := shared.InitServerSpan(ctx, ServiceName)
	defer span.End()

	latency, _ := strconv.ParseInt(os.Getenv("CONSTANT_LATENCY"), 10, 32)
	time.Sleep(time.Duration(latency) * time.Millisecond)

	span.SetStatus(codes.Ok, "success")

	f := &Fish.Single{
		Id:             shared.GenerateRandomUUID(),
		RandomMetadata: shared.GenerateFakeMetadata(),
	}

	return f, nil
}

// --------------- HTTP Methods ---------------

func ProduceFish(w http.ResponseWriter, r *http.Request) {
	ctx, span := shared.InitServerSpan(context.Background(), ServiceName)
	defer span.End()

	r.WithContext(ctx)
	w.Header().Set("Content-Type", "application/json")

	latency, _ := strconv.ParseInt(os.Getenv("CONSTANT_LATENCY"), 10, 32)
	time.Sleep(time.Duration(latency) * time.Millisecond)

	d := shared.FishHTTP{
		FishId:             shared.GenerateRandomUUID(),
		FishRandomMetadata: shared.GenerateFakeMetadata(),
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
		Fish.RegisterFishServer(grpcServer, NewFishServer(logger))

		go func() {
			if err := grpcServer.Serve(listener); err != nil {
				logger.Fatalf("could not start grpc server: %v", err)
			}
		}()

		shared.MonitorShutdownSignal()
	} else {
		logger.Info("Using HTTP")
		mux := http.NewServeMux()
		mux.HandleFunc("/", ProduceFish)

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
