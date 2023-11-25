package shared

import (
	"context"
	"crypto/rand"
	"fmt"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/microsoft/go-mssqldb"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"log"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"syscall"
)

// --------------- Application-Related Operations ---------------

func MonitorShutdownSignal() {
	println("Monitoring shutdown signal")
	signalChan := make(chan os.Signal, 1)
	signal.Notify(
		signalChan,
		syscall.SIGTERM, // https://cloud.google.com/blog/topics/developers-practitioners/graceful-shutdowns-cloud-run-deep-dive
		syscall.SIGHUP,  // kill -SIGHUP
		syscall.SIGINT,  // kill -SIGINT or Ctrl+c
		syscall.SIGQUIT, // kill -SIGQUIT
	)

	<-signalChan
	log.Printf("os.Interrupt - shutting down...\n")

	// terminate after second signal before callback is done
	go func() {
		<-signalChan
		log.Printf("os.Kill - terminating...\n")
	}()

	// PERFORM GRACEFUL SHUTDOWN HERE
	os.Exit(0)
}

func InitSugaredLogger() *otelzap.SugaredLogger {
	logger, err := zap.NewProduction()
	if err != nil {
		println("failed to init production logger; exiting...")
		os.Exit(1)
	}

	sLogger := otelzap.New(logger).Sugar()

	return sLogger
}

func ConfigureRuntime() {
	//nuCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(2)
	fmt.Printf("RUNNING WITH %d CPU\n", 4)
}

func LoadEnvFile(additionalEnv string) {
	err := godotenv.Load("./shared_env", additionalEnv)
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

// --------------- Observability-Related Operations ---------------

func SetGRPCHeader(ctx *context.Context) {
	header := metadata.Pairs("Content-Type", "application/grpc")
	err := grpc.SetHeader(*ctx, header)
	if err != nil {
		log.Printf("failed to set header to gRPC Request: %v", err)
	}
}

func getCallFuncName() string {
	spanName := ""
	pc, _, _, ok := runtime.Caller(2)
	details := runtime.FuncForPC(pc)
	if ok && details != nil {
		callNames := strings.Split(details.Name(), ".")
		spanName = callNames[len(callNames)-1]

		fmt.Printf("called from %s\n", spanName)
	}

	return spanName
}

func InitServerSpan(ctx context.Context, name string) (context.Context, trace.Span) {
	return otel.Tracer(name).Start(ctx, getCallFuncName(), trace.WithSpanKind(trace.SpanKindServer))
}

func InitInternalSpan(ctx context.Context) (context.Context, trace.Span) {
	return otel.Tracer("").Start(ctx, getCallFuncName())
}

// --------------- Prometheus Metrics ---------------

func InitPrometheusRequestCountMetrics() *prometheus.CounterVec {
	var (
		RequestCount = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "request_to_node_count",
				Help: "Count the number of incoming requests to the particular node",
			},
			[]string{"service_name", "node_name"}, // Labels for the metric, if any
		)
	)

	//var (
	//	OutgoingRequestCount = prometheus.NewCounterVec(
	//		prometheus.CounterOpts{
	//			Name: "outgoing_request_to_node_count",
	//			Help: "Count the number of outgoing requests to the particular node",
	//		},
	//		[]string{"serviceName", "nodeName"},
	//	)
	//)

	return RequestCount
}

func IncreaseRequestCount(nodeName string, serviceName string, vec *prometheus.CounterVec) {

}

// --------------- gRPC Related ---------------

func InitGrpcClientConn(targetAddr string) *grpc.ClientConn {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	var conn *grpc.ClientConn
	var err error

	conn, err = grpc.Dial(targetAddr, opts...)
	if err != nil {
		panic(err)
	}

	return conn
}

// --------------- Shared Data Structure ---------------

func GenerateRandomUUID() string {
	return uuid.New().String()
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateFakeMetadataString(ctx context.Context, size string) string {
	InitInternalSpan(ctx)

	sizeInByte, err := strconv.ParseUint(strings.TrimRight(size, "bkm"), 10, 64)
	if err != nil {
		sizeInByte = 1
	}

	if strings.HasSuffix(size, "k") {
		sizeInByte *= 1024
	} else if strings.HasSuffix(size, "m") {
		sizeInByte *= 1024 * 1024
	}

	bytes := make([]byte, sizeInByte)
	rand.Read(bytes)

	for i := range bytes {
		bytes[i] = charset[bytes[i]%byte(len(charset))]
	}

	// You can convert the random bytes to a string using base64 encoding or any other method you prefer
	randomString := string(bytes)

	return randomString
}

type SingleBasicType struct {
	Id             string
	RandomMetadata string
}

type BasicTypeHTTPResponse struct {
	Quantity uint64
	Type     string
	Items    []SingleBasicType
}

type SingleFour struct {
	Id             string
	RandomMetadata string
	GrainId        string
}

type FlourHTTPResponse struct {
	Quantity uint64
	Type     string
	Items    []SingleFour
}

type SingleLog struct {
	Id             string
	RandomMetadata string
	TreeId         string
}

type LogHTTPResponse struct {
	Quantity uint64
	Type     string
	Items    []SingleLog
}

type SingleBoard struct {
	Id             string
	RandomMetadata string
	LogId          string
}

type BoardHTTPResponse struct {
	Quantity uint64
	Type     string
	Items    []SingleBoard
}
