package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kubernauts/tk8/pkg/api"
	"github.com/rs/cors"
)

func main() {
	fmt.Println("Hello World")
	r := mux.NewRouter()
	r.HandleFunc("/", api.DemoHandler).Methods("GET")
	r.HandleFunc("/cluster/init", api.InfraOnlyHandler).Methods("POST")
	r.HandleFunc("/cluster/create", api.CreateHandler).Methods("POST")
	log.Fatal(http.ListenAndServe(":8081", cors.Default().Handler(r)))
}
