package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"service-date/service"
	"service-date/transport"
	"syscall"
)

func main() {
	var (
		httpAddr = flag.String("http", ":8080", "http listen address")
	)
	flag.Parse()
	ctx := context.Background()
	// our napodate service
	srv := service.NewService()
	errChan := make(chan error)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	// mapping endpoints
	endpoints := transport.Endpoints{
		GetEndpoint:      transport.MakeGetEndpoint(srv),
		StatusEndpoint:   transport.MakeStatusEndpoint(srv),
		ValidateEndpoint: transport.MakeValidateEndpoint(srv),
	}

	// HTTP transport
	go func() {
		log.Println("service-date is listening on port:", *httpAddr)
		handler := transport.NewHTTPServer(ctx, endpoints)
		errChan <- http.ListenAndServe(*httpAddr, handler)
	}()

	log.Fatalln(<-errChan)
}