package shared

import (
	Bread "NetMARKS/services/bread/proto"
	Fish "NetMARKS/services/fish/proto"
	Flour "NetMARKS/services/flour/proto"
	Grain "NetMARKS/services/grain/proto"
	Meat "NetMARKS/services/meat/proto"
	Water "NetMARKS/services/water/proto"
	"context"
	"crypto/rand"
	"fmt"
	"github.com/google/uuid"
	_ "github.com/microsoft/go-mssqldb"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

// --------------- Application-Related Operations ---------------

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
	return prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "request_to_node_count",
			Help: "Count the number of incoming requests to the particular node",
		},
		[]string{"service_name", "node_name"}, // Labels for the metric, if any
	)
}

// --------------- gRPC Related ---------------

type GRPCResponse struct {
	Type string
	Body string
	Err  error
}

func ConcurrentGRPCWater(ctx context.Context, client Water.WaterClient, wg *sync.WaitGroup, ch chan<- GRPCResponse) {
	t := "water"
	defer wg.Done()

	produce, err := client.Produce(ctx, &Water.Request{
		Quantity:     1,
		ResponseSize: "1",
	})
	if err != nil {
		ch <- GRPCResponse{Type: t, Err: err}
		return
	}

	ch <- GRPCResponse{Type: t, Body: produce.Items[0].Id}
}

func ConcurrentGRPCGrain(ctx context.Context, client Grain.GrainClient, wg *sync.WaitGroup, ch chan<- GRPCResponse) {
	t := "grain"
	defer wg.Done()

	produce, err := client.Produce(ctx, &Grain.Request{
		Quantity:     1,
		ResponseSize: "1",
	})
	if err != nil {
		ch <- GRPCResponse{Type: t, Err: err}
		return
	}

	ch <- GRPCResponse{Type: t, Body: produce.Items[0].Id}
}

func ConcurrentGRPCFlour(ctx context.Context, client Flour.FlourClient, wg *sync.WaitGroup, ch chan<- GRPCResponse) {
	t := "flour"
	defer wg.Done()

	produce, err := client.Produce(ctx, &Flour.Request{
		Quantity:     1,
		ResponseSize: "1",
	})
	if err != nil {
		ch <- GRPCResponse{Type: t, Err: err}
		return
	}

	ch <- GRPCResponse{Type: t, Body: produce.Items[0].Id}
}

func ConcurrentGRPCBread(ctx context.Context, client Bread.BreadClient, wg *sync.WaitGroup, ch chan<- GRPCResponse) {
	t := "bread"
	defer wg.Done()

	produce, err := client.Produce(ctx, &Bread.Request{
		Quantity:     1,
		ResponseSize: "1",
	})
	if err != nil {
		ch <- GRPCResponse{Type: t, Err: err}
		return
	}

	ch <- GRPCResponse{Type: t, Body: produce.Items[0].Id}
}

func ConcurrentGRPCFish(ctx context.Context, client Fish.FishClient, wg *sync.WaitGroup, ch chan<- GRPCResponse) {
	t := "fish"
	defer wg.Done()

	produce, err := client.Produce(ctx, &Fish.Request{
		Quantity:     1,
		ResponseSize: "1",
	})
	if err != nil {
		ch <- GRPCResponse{Type: t, Err: err}
		return
	}

	ch <- GRPCResponse{Type: t, Body: produce.Items[0].Id}
}

func ConcurrentGRPCMeat(ctx context.Context, client Meat.MeatClient, wg *sync.WaitGroup, ch chan<- GRPCResponse) {
	t := "meat"
	defer wg.Done()

	produce, err := client.Produce(ctx, &Meat.Request{
		Quantity:     1,
		ResponseSize: "1",
	})
	if err != nil {
		ch <- GRPCResponse{Type: t, Err: err}
		return
	}

	ch <- GRPCResponse{Type: t, Body: produce.Items[0].Id}
}

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

// --------------- HTTP Related ---------------

type HTTPResponse struct {
	Type string
	Body []byte
	Err  error
}

func ConcurrentHTTPRequest(url string, productType string, wg *sync.WaitGroup, ch chan<- HTTPResponse) {
	defer wg.Done()

	resp, err := http.Get(url)
	if err != nil {
		ch <- HTTPResponse{Type: productType, Err: err}
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		ch <- HTTPResponse{Type: productType, Err: err}
		return
	}

	ch <- HTTPResponse{Type: productType, Body: body}
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
	Id             string `json:"id,omitempty"`
	RandomMetadata string `json:"randomMetadata,omitempty"`
}

type BasicTypeHTTPResponse struct {
	Quantity uint64            `json:"quantity,omitempty"`
	Type     string            `json:"type,omitempty"`
	Items    []SingleBasicType `json:"items,omitempty"`
}

type SingleFour struct {
	Id             string `json:"id,omitempty"`
	RandomMetadata string `json:"randomMetadata,omitempty"`
	GrainId        string `json:"grainId,omitempty"`
}

type FlourHTTPResponse struct {
	Quantity uint64       `json:"quantity,omitempty"`
	Type     string       `json:"type,omitempty"`
	Items    []SingleFour `json:"items,omitempty"`
}

type SingleLog struct {
	Id             string `json:"id,omitempty"`
	RandomMetadata string `json:"randomMetadata,omitempty"`
	TreeId         string `json:"treeId,omitempty"`
}

type LogHTTPResponse struct {
	Quantity uint64      `json:"quantity,omitempty"`
	Type     string      `json:"type,omitempty"`
	Items    []SingleLog `json:"items,omitempty"`
}

type SingleBoard struct {
	Id             string `json:"id,omitempty"`
	RandomMetadata string `json:"randomMetadata,omitempty"`
	LogId          string `json:"logId,omitempty"`
}

type BoardHTTPResponse struct {
	Quantity uint64        `json:"quantity,omitempty"`
	Type     string        `json:"type,omitempty"`
	Items    []SingleBoard `json:"items,omitempty"`
}

type SingleBeer struct {
	Id             string `json:"id,omitempty"`
	RandomMetadata string `json:"randomMetadata,omitempty"`
	GrainId        string `json:"grainId,omitempty"`
	WaterId        string `json:"waterId,omitempty"`
}

type BeerHTTPResponse struct {
	Quantity uint64       `json:"quantity,omitempty"`
	Type     string       `json:"type,omitempty"`
	Items    []SingleBeer `json:"items,omitempty"`
}

type SinglePig struct {
	Id             string `json:"id,omitempty"`
	RandomMetadata string `json:"randomMetadata,omitempty"`
	GrainId        string `json:"grainId,omitempty"`
	WaterId        string `json:"waterId,omitempty"`
}

type PigHTTPResponse struct {
	Quantity uint64      `json:"quantity,omitempty"`
	Type     string      `json:"type,omitempty"`
	Items    []SinglePig `json:"items,omitempty"`
}

type SingleBoat struct {
	Id             string `json:"id,omitempty"`
	RandomMetadata string `json:"randomMetadata,omitempty"`
	BoardId        string `json:"boardId,omitempty"`
}

type BoatHTTPResponse struct {
	Quantity uint64       `json:"quantity,omitempty"`
	Type     string       `json:"type,omitempty"`
	Items    []SingleBoat `json:"items,omitempty"`
}

type SingleMeat struct {
	Id             string `json:"id,omitempty"`
	RandomMetadata string `json:"randomMetadata,omitempty"`
	PigId          string `json:"pigId,omitempty"`
}

type MeatHTTPResponse struct {
	Quantity uint64       `json:"quantity,omitempty"`
	Type     string       `json:"type,omitempty"`
	Items    []SingleMeat `json:"items,omitempty"`
}

type SingleBread struct {
	Id             string `json:"id,omitempty"`
	RandomMetadata string `json:"randomMetadata,omitempty"`
	WaterId        string `json:"waterId,omitempty"`
	FlourId        string `json:"flourId,omitempty"`
}

type BreadHTTPResponse struct {
	Quantity uint64        `json:"quantity,omitempty"`
	Type     string        `json:"type,omitempty"`
	Items    []SingleBread `json:"items,omitempty"`
}

type SingleCoal struct {
	Id             string `json:"id,omitempty"`
	RandomMetadata string `json:"randomMetadata,omitempty"`
	FishId         string `json:"fishId,omitempty"`
	BreadId        string `json:"breadId,omitempty"`
	MeatId         string `json:"meatId,omitempty"`
}

type CoalHTTPResponse struct {
	Quantity uint64       `json:"quantity,omitempty"`
	Type     string       `json:"type,omitempty"`
	Items    []SingleCoal `json:"items,omitempty"`
}
