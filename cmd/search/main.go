package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ris.com/internal"
	"ris.com/internal/controller"
	"ris.com/internal/gateway"
	"ris.com/internal/handlers"
	"ris.com/internal/repository"
)

func getServer() *http.Server {
	db, err := internal.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to db: %v", err)
	}

	client := getClient()

	ir := repository.NewImageRepository(db)
	eg := gateway.NewEmbeddingGateway(client)
	sc := controller.NewSearchController(ir, eg)
	ssh := handlers.NewSearchSimilarHandler(sc)

	mux := http.NewServeMux()
	mux.Handle("/api/v1/searchSimilar", ssh)

	server := &http.Server{
		Addr: ":8080",
		Handler: mux,
	}
	return server
}

func getClient() *http.Client {
	transportConfig := &http.Transport{
		MaxIdleConns: 100,
		MaxIdleConnsPerHost: 100,
		IdleConnTimeout: 90 * time.Second,
		DisableCompression: false,
		ExpectContinueTimeout: 1 * time.Second,
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: transportConfig,
	}
	return client
}

func main() {
	server := getServer()

	_ = context.Background()

	// sc.SearchSimilar(ctx, nil, 5, 0)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Println("Starting server on port 8080")
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP Server error: %v", err)
		}
		log.Println("Stopped server")
	}()

	<-stop // blocks until stop os.Interrupt or syscall.SIGTERM is sent to process
	log.Println("Attempting server shutdown...")

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Forced shutdown: %v", err)
	}
}

