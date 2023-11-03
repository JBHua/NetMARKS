package main

import (
	"VerbosityMonoRepo/shared"
	"fmt"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"net/http"
	"os"
)

func UploadMessage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %s", os.Getenv("CITY"))
}

type FishServer struct {
	VUser.UnimplementedVerbosityUserServer
	repo   *UserRepo
	logger *otelzap.SugaredLogger
}

func ProduceFishHTTP(w http.ResponseWriter, r *http.Request) {

}
func ProduceFishGRPC() {

}

func main() {
	shared.ConfigureRuntime()

	fmt.Println("Hello???")

	mux := http.NewServeMux()
	mux.HandleFunc("/", UploadMessage)

	// Start HTTP Server
	port := os.Getenv("PORT")
	fmt.Println(port)
	err := http.ListenAndServe(":"+port, mux)
	if err != nil {
		panic(err)
	}
	fmt.Printf("service running on port %s\n", port)
}
