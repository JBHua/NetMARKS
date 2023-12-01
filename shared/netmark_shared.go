package shared

import (
	Board "NetMARKS/services/board/proto"
	Bread "NetMARKS/services/bread/proto"
	Coal "NetMARKS/services/coal/proto"
	Fish "NetMARKS/services/fish/proto"
	Flour "NetMARKS/services/flour/proto"
	Gold "NetMARKS/services/gold/proto"
	Grain "NetMARKS/services/grain/proto"
	Iron "NetMARKS/services/iron/proto"
	Ironore "NetMARKS/services/ironore/proto"
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
	runtime.GOMAXPROCS(4)
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

func InitPrometheusRequestCountMetrics() (*prometheus.CounterVec, *prometheus.CounterVec) {
	requestCount := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "sub_request_count_per_service",
			Help: "Count the number of sub-request for a specific service",
		},
		[]string{"service_name", "node_name", "upstream_node_name", "original_request_service"},
	)

	interNodeRequestCount := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "inter_node_sub_request_count_per_service",
			Help: "Count the number of sub-request coming from another node for a service",
		},
		[]string{"service_name", "node_name", "upstream_node_name", "original_request_service"},
	)

	return requestCount, interNodeRequestCount
}

func UpdateRequestMetrics(requestCounter, interNodeRequestCounter *prometheus.CounterVec, originalRequestService, serviceName, upstreamNodeName, nodeName string) {
	fmt.Printf("Request Metrics for: %s on Node: %s. Upstreaming request is from: %s on Node: %s\n", serviceName, nodeName, originalRequestService, upstreamNodeName)
	if originalRequestService != serviceName {
		IncreaseRequestCount(requestCounter, originalRequestService, serviceName, upstreamNodeName, nodeName)
	}

	if upstreamNodeName != nodeName {
		IncreaseRequestCount(interNodeRequestCounter, originalRequestService, serviceName, upstreamNodeName, nodeName)
	}
}

func IncreaseRequestCount(counter *prometheus.CounterVec, originalRequestService, serviceName, upstreamNodeName, nodeName string) {
	counter.With(prometheus.Labels{
		"service_name":             serviceName,
		"node_name":                nodeName,
		"upstream_node_name":       upstreamNodeName,
		"original_request_service": originalRequestService,
	}).Inc()
}

// --------------- gRPC Related ---------------

type GRPCResponse struct {
	Type string
	Body string
	Err  error
}

func ConcurrentGRPCWater(ctx context.Context, size string, client Water.WaterClient, wg *sync.WaitGroup, ch chan<- GRPCResponse) {
	t := "water"
	defer wg.Done()

	produce, err := client.Produce(ctx, &Water.Request{
		Quantity:     1,
		ResponseSize: size,
	})
	if err != nil {
		ch <- GRPCResponse{Type: t, Err: err}
		return
	}

	ch <- GRPCResponse{Type: t, Body: produce.Items[0].Id}
}

func ConcurrentGRPCGrain(ctx context.Context, size string, client Grain.GrainClient, wg *sync.WaitGroup, ch chan<- GRPCResponse) {
	t := "grain"
	defer wg.Done()

	produce, err := client.Produce(ctx, &Grain.Request{
		Quantity:     1,
		ResponseSize: size,
	})
	if err != nil {
		ch <- GRPCResponse{Type: t, Err: err}
		return
	}

	ch <- GRPCResponse{Type: t, Body: produce.Items[0].Id}
}

func ConcurrentGRPCFlour(ctx context.Context, size string, client Flour.FlourClient, wg *sync.WaitGroup, ch chan<- GRPCResponse) {
	t := "flour"
	defer wg.Done()

	produce, err := client.Produce(ctx, &Flour.Request{
		Quantity:     1,
		ResponseSize: size,
	})
	if err != nil {
		ch <- GRPCResponse{Type: t, Err: err}
		return
	}

	ch <- GRPCResponse{Type: t, Body: produce.Items[0].Id}
}

func ConcurrentGRPCBread(ctx context.Context, size string, client Bread.BreadClient, wg *sync.WaitGroup, ch chan<- GRPCResponse) {
	t := "bread"
	defer wg.Done()

	produce, err := client.Produce(ctx, &Bread.Request{
		Quantity:     1,
		ResponseSize: size,
	})
	if err != nil {
		ch <- GRPCResponse{Type: t, Err: err}
		return
	}

	ch <- GRPCResponse{Type: t, Body: produce.Items[0].Id}
}

func ConcurrentGRPCFish(ctx context.Context, size string, client Fish.FishClient, wg *sync.WaitGroup, ch chan<- GRPCResponse) {
	t := "fish"
	defer wg.Done()

	produce, err := client.Produce(ctx, &Fish.Request{
		Quantity:     1,
		ResponseSize: size,
	})
	if err != nil {
		ch <- GRPCResponse{Type: t, Err: err}
		return
	}

	ch <- GRPCResponse{Type: t, Body: produce.Items[0].Id}
}

func ConcurrentGRPCMeat(ctx context.Context, size string, client Meat.MeatClient, wg *sync.WaitGroup, ch chan<- GRPCResponse) {
	t := "meat"
	defer wg.Done()

	produce, err := client.Produce(ctx, &Meat.Request{
		Quantity:     1,
		ResponseSize: size,
	})
	if err != nil {
		ch <- GRPCResponse{Type: t, Err: err}
		return
	}

	ch <- GRPCResponse{Type: t, Body: produce.Items[0].Id}
}

func ConcurrentGRPCCoal(ctx context.Context, size string, client Coal.CoalClient, wg *sync.WaitGroup, ch chan<- GRPCResponse) {
	t := "coal"
	defer wg.Done()

	produce, err := client.Produce(ctx, &Coal.Request{
		Quantity:     1,
		ResponseSize: size,
	})
	if err != nil {
		ch <- GRPCResponse{Type: t, Err: err}
		return
	}

	ch <- GRPCResponse{Type: t, Body: produce.Items[0].Id}
}

func ConcurrentGRPCGold(ctx context.Context, size string, client Gold.GoldClient, wg *sync.WaitGroup, ch chan<- GRPCResponse) {
	t := "gold"
	defer wg.Done()

	produce, err := client.Produce(ctx, &Gold.Request{
		Quantity:     1,
		ResponseSize: size,
	})
	if err != nil {
		ch <- GRPCResponse{Type: t, Err: err}
		return
	}

	ch <- GRPCResponse{Type: t, Body: produce.Items[0].Id}
}

func ConcurrentGRPCIronore(ctx context.Context, size string, client Ironore.IronoreClient, wg *sync.WaitGroup, ch chan<- GRPCResponse) {
	t := "gold"
	defer wg.Done()

	produce, err := client.Produce(ctx, &Ironore.Request{
		Quantity:     1,
		ResponseSize: size,
	})
	if err != nil {
		ch <- GRPCResponse{Type: t, Err: err}
		return
	}

	ch <- GRPCResponse{Type: t, Body: produce.Items[0].Id}
}

func ConcurrentGRPCBoard(ctx context.Context, size string, client Board.BoardClient, wg *sync.WaitGroup, ch chan<- GRPCResponse) {
	t := "gold"
	defer wg.Done()

	produce, err := client.Produce(ctx, &Board.Request{
		Quantity:     1,
		ResponseSize: size,
	})
	if err != nil {
		ch <- GRPCResponse{Type: t, Err: err}
		return
	}

	ch <- GRPCResponse{Type: t, Body: produce.Items[0].Id}
}

func ConcurrentGRPCIron(ctx context.Context, size string, client Iron.IronClient, wg *sync.WaitGroup, ch chan<- GRPCResponse) {
	t := "gold"
	defer wg.Done()

	produce, err := client.Produce(ctx, &Iron.Request{
		Quantity:     1,
		ResponseSize: size,
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

func ConcurrentHTTPRequest(url, productType, nodeName, requestID, originalRequestService string, wg *sync.WaitGroup, ch chan<- HTTPResponse) {
	defer wg.Done()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		ch <- HTTPResponse{Type: productType, Err: err}
		return
	}

	req.Header.Set("upstream-node-name", nodeName)
	req.Header.Set("original-request-service", originalRequestService)
	req.Header.Set("request-id", requestID)

	resp, err := http.DefaultClient.Do(req)
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

func ExtractUpstreamRequestID(header http.Header, serviceName string, nodeName string) (string, string, string) {
	var requestID string
	var originalRequestService string
	var upstreamNodeName string

	id := header.Get("request-id")
	if len(id) > 0 {
		fmt.Printf("Getting Existing Request ID: %s\n", id)
		requestID = id
	} else {
		fmt.Printf("No upstream request. Generating new id....")
		requestID = GenerateRandomUUID()
	}

	originalService := header.Get("original-request-service")
	if len(originalService) > 0 {
		originalRequestService = originalService
	} else {
		originalRequestService = serviceName
	}

	upstream := header.Get("upstream-node-name")
	if len(upstream) > 0 {
		upstreamNodeName = upstream
	} else {
		upstreamNodeName = nodeName
	}

	return requestID, originalRequestService, upstreamNodeName
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

type SingleCoin struct {
	Id             string `json:"id,omitempty"`
	RandomMetadata string `json:"randomMetadata,omitempty"`
	CoalId         string `json:"coalId,omitempty"`
	GoldId         string `json:"goldId,omitempty"`
}

type CoinHTTPResponse struct {
	Quantity uint64       `json:"quantity,omitempty"`
	Type     string       `json:"type,omitempty"`
	Items    []SingleCoin `json:"items,omitempty"`
}

type SingleIron struct {
	Id             string `json:"id,omitempty"`
	RandomMetadata string `json:"randomMetadata,omitempty"`
	CoalId         string `json:"coalId,omitempty"`
	IronoreId      string `json:"ironoreId,omitempty"`
}

type IronHTTPResponse struct {
	Quantity uint64       `json:"quantity,omitempty"`
	Type     string       `json:"type,omitempty"`
	Items    []SingleIron `json:"items,omitempty"`
}

type SingleTool struct {
	Id             string `json:"id,omitempty"`
	RandomMetadata string `json:"randomMetadata,omitempty"`
	BoardId        string `json:"boardId,omitempty"`
	IronoreId      string `json:"ironoreId,omitempty"`
}

type ToolsHTTPResponse struct {
	Quantity uint64       `json:"quantity,omitempty"`
	Type     string       `json:"type,omitempty"`
	Items    []SingleTool `json:"items,omitempty"`
}

type SingleSword struct {
	Id             string `json:"id,omitempty"`
	RandomMetadata string `json:"randomMetadata,omitempty"`
	CoalId         string `json:"coalId,omitempty"`
	IronId         string `json:"ironId,omitempty"`
}

type SwordHTTPResponse struct {
	Quantity uint64        `json:"quantity,omitempty"`
	Type     string        `json:"type,omitempty"`
	Items    []SingleSword `json:"items,omitempty"`
}
