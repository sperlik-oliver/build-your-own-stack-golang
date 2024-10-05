package main

import (
	"build-your-own/load-balancer/config"
	"build-your-own/load-balancer/network"
	"fmt"
	"log"
	"net/http"
)

// Entrypoint
func main() {
	config.LoadConfig()
	go network.Healthcheck()

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", config.GetConfig().Port),
		Handler: http.HandlerFunc(network.Lb),
	}

	log.Printf("Load balancer started at :%d", config.GetConfig().Port)
	server.ListenAndServe()
}
