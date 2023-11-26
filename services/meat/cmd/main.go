package main

import (
	Meat "NetMARKS/services/meat/proto"
	Pig "NetMARKS/services/pig/proto"
	"NetMARKS/shared"
	"context"
	"encoding/json"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/soheilhy/cmux"
	"go.opentelemetry.io/otel/codes"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"
)

const ServiceName = "Meat"
const ServicePort = "8080"
const PigServiceAddr = "netmarks-pig.default.svc.cluster.local:8080"

var NodeName = os.Getenv("K8S_NODE_NAME")
var RequestCount = shared.InitPrometheusRequestCountMetrics()

// --------------- gRPC Methods ---------------

type MeatServer struct {
	Meat.UnimplementedMeatServer
	pigClient Pig.PigClient
}

func NewMeatServer(p Pig.PigClient) *MeatServer {
	return &MeatServer{
		pigClient: p,
	}
}

func newGRPCServer(lis net.Listener) error {
	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	pigClient := Pig.NewPigClient(shared.InitGrpcClientConn(PigServiceAddr))
	Meat.RegisterMeatServer(grpcServer, NewMeatServer(pigClient))

	return grpcServer.Serve(lis)
}

func (s *MeatServer) Produce(ctx context.Context, req *Meat.Request) (*Meat.Response, error) {
	shared.SetGRPCHeader(&ctx)
	ctx, span := shared.InitServerSpan(ctx, ServiceName)
	defer span.End()
	defer RequestCount.With(prometheus.Labels{
		"service_name": ServiceName,
		"node_name":    NodeName,
	}).Inc()

	latency, _ := strconv.ParseInt(os.Getenv("CONSTANT_LATENCY"), 10, 32)

	r := Meat.Response{}
	for i := uint64(0); i < req.Quantity; i++ {
		r.Quantity += 1

		produce, err := s.pigClient.Produce(ctx, &Pig.Request{
			Quantity:     1,
			ResponseSize: "1",
		})
		if err != nil {
			return nil, err
		}

		r.Items = append(r.Items, &Meat.Single{
			Id:             shared.GenerateRandomUUID(),
			RandomMetadata: shared.GenerateFakeMetadataString(ctx, req.ResponseSize),
			PigId:          produce.Items[0].Id,
		})

		time.Sleep(time.Duration(latency) * time.Millisecond)
	}

	span.SetStatus(codes.Ok, "success")
	return &r, nil
}

// --------------- HTTP Methods ---------------

func newHTTPServer(lis net.Listener) error {
	mux := http.NewServeMux()

	mux.HandleFunc("/", Produce)
	mux.Handle("/metrics", promhttp.Handler())

	s := &http.Server{Handler: mux}
	return s.Serve(lis)
}

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

	latency, _ := strconv.ParseInt(os.Getenv("CONSTANT_LATENCY"), 10, 32)

	response := shared.MeatHTTPResponse{
		Type: ServiceName,
	}
	for i := uint64(0); i < quantity; i++ {
		response.Quantity += 1

		getRes, err := http.Get("http://" + PigServiceAddr + "?quantity=1&response_size=1")
		if err != nil {
			w.Write([]byte(err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer getRes.Body.Close()

		var pig shared.PigHTTPResponse
		err = json.NewDecoder(getRes.Body).Decode(&pig)
		if err != nil {
			w.Write([]byte(err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		response.Items = append(response.Items, shared.SingleMeat{
			Id:             shared.GenerateRandomUUID(),
			RandomMetadata: shared.GenerateFakeMetadataString(ctx, r.URL.Query().Get("response_size")),
			PigId:          pig.Items[0].Id,
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

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", ServicePort))
	if err != nil {
		logger.Fatalf("could not attach listener to port: %v. %v", ServicePort, err)
	}

	mux := cmux.New(listener)
	grpcListener := mux.Match(cmux.HTTP2HeaderField("content-type", "application/grpc"))
	httpListener := mux.Match(cmux.HTTP1Fast())

	// Use an error group to start all of them
	g := errgroup.Group{}
	g.Go(func() error { return newGRPCServer(grpcListener) })
	g.Go(func() error { return newHTTPServer(httpListener) })
	g.Go(func() error { return mux.Serve() })

	log.Println("run server:", g.Wait())
}
