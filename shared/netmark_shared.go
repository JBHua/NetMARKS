package shared

import (
	"context"
	"crypto/rand"
	"fmt"
	"github.com/google/uuid"
	"github.com/johnsiilver/getcert"
	"github.com/joho/godotenv"
	_ "github.com/microsoft/go-mssqldb"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"log"
	"os"
	"os/signal"
	"runtime"
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
	nuCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(nuCPU)
	fmt.Printf("RUNNING WITH %d CPU\n", nuCPU)
}

func LoadEnvFile(additionalEnv string) {
	err := godotenv.Load("./shared_env", additionalEnv)
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

// --------------- Observability-Related Operations ---------------

func InitTracerProvider(debug bool, service string) (*sdktrace.TracerProvider, error) {
	var secureOption otlptracegrpc.Option
	var addr string

	if !debug {
		addr = os.Getenv("OTLP_COLLECTOR_ENDPOINT")
		// HTTP2 requires us to use https, thus, no more otlptracegrpc.WithInsecure()
		// https://kennethjenkins.net/posts/go-nginx-grpc/
		tlsCert, _, _ := getcert.FromTLSServer(addr, true)
		secureOption = otlptracegrpc.WithTLSCredentials(credentials.NewServerTLSFromCert(&tlsCert))
	} else {
		addr = os.Getenv("OTLP_COLLECTOR_ENDPOINT_DEBUG")
		secureOption = otlptracegrpc.WithInsecure()
	}

	exporter, err := otlptrace.New(
		context.Background(),
		otlptracegrpc.NewClient(
			secureOption,
			otlptracegrpc.WithEndpoint(addr),
		),
	)
	if err != nil {
		log.Fatalf("failed to start a trace exporter ")
		return nil, err
	}

	appResource, _ := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(os.Getenv(strings.ToUpper(service)+"_SERVICE_NAME")),
			semconv.ServiceVersionKey.String(os.Getenv(strings.ToUpper(service)+"_SERVICE_VERSION")),
		),
	)

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(appResource),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return tp, nil
}

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

// --------------- Shared Data Structure ---------------

func InitGrpcClientConn(targetAddr string) *grpc.ClientConn {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	var conn *grpc.ClientConn
	var err error

	conn, err = grpc.Dial("localhost:"+targetAddr, opts...)
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

func GenerateFakeMetadataInKB(ctx context.Context, sizeInKB uint64) string {
	InitInternalSpan(ctx)

	bytes := make([]byte, sizeInKB*1024)
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

type WaterHTTP struct {
	Id             string
	RandomMetadata string
}
