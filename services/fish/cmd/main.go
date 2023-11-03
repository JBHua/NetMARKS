package main

import (
	"NetMARKS/shared"
	"fmt"
	Fish "github.com/JBHua/NetMARKS/services/fish/proto"
	"github.com/joho/godotenv"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
)

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

func ProduceFishGRPC() {

}

// --------------- HTTP Methods ---------------

func UploadMessage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %s", os.Getenv("CITY"))
}

func ProduceFishHTTP(w http.ResponseWriter, r *http.Request) {

}

func main() {
	err := godotenv.Load("../../shared/.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	shared.ConfigureRuntime()
	logger := shared.InitSugaredLogger()

	useGRPC, _ := strconv.ParseBool(os.Getenv("USE_GRPC"))
	if useGRPC {
		println("Using GRPC")
		listener, err := net.Listen("tcp", fmt.Sprintf(":%s", os.Getenv(ServicePortEnv)))
		if err != nil {
			log.Fatalf("could not attach listener to port: %v. %v", os.Getenv(ServicePortEnv), err)
		}
		fmt.Printf("Running at %s\n", os.Getenv(ServicePortEnv))

		grpcServer := grpc.NewServer()
		Fish.RegisterFishServer(grpcServer, NewFishServer(logger))

		go func() {
			if err := grpcServer.Serve(listener); err != nil {
				logger.Fatalf("could not start grpc server: %v", err)
			}
		}()

		shared.MonitorShutdownSignal()
	} else {
		println("Using HTTP")
		mux := http.NewServeMux()
		mux.HandleFunc("/", UploadMessage)

		// Start HTTP Server
		port := os.Getenv(ServicePortEnv)
		fmt.Println(port)
		err = http.ListenAndServe(":"+port, mux)

		if err != nil {
			panic(err)
		}
		fmt.Printf("service running on port %s\n", port)
	}
}
