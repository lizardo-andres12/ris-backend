package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"ris.com/internal/controller"
	"ris.com/internal/database"
	"ris.com/internal/gateway"
	"ris.com/internal/handlers"
	"ris.com/internal/http/client"
	"ris.com/internal/otel"
	"ris.com/internal/repository"
)

func main() {
	var shutdownHooks []func(context.Context) error
	startupCtx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	resource, err := otel.NewResource()
	if err != nil {
		panic(err)
	}
	log.Println("hello")

	logExporter, err := otel.NewLogExporter()
	if err != nil {
		panic(err)
	}
	log.Println("hello")

	metricExporter, err := otel.NewMetricExporter()
	if err != nil {
		panic(err)
	}
	log.Println("nye")

	traceExporter, err := otel.NewTraceExporter()
	if err != nil {
		panic(err)
	}
	log.Println("hello")

	db, err := database.NewDB(startupCtx)
	if err != nil {
		panic(err)
	}
	log.Println("hello")

	log.Println("Creating logger provider...")
	lp := otel.NewLoggerProvider(resource, logExporter)
	shutdownHooks = append(shutdownHooks, lp.Shutdown)
	log.Println("bye")

	log.Println("Creating metric provider...")
	mp := otel.NewMetricProvider(resource, metricExporter)
	shutdownHooks = append(shutdownHooks, mp.Shutdown)
	log.Println("hello")

	log.Println("Creating trace provider...")
	tp := otel.NewTraceProvider(resource, traceExporter)
	shutdownHooks = append(shutdownHooks, tp.Shutdown)
	log.Println("hello")

	httpClient := client.NewHTTPClient(client.NewHTTPClientTransport())

	embeddingGateway := gateway.NewEmbeddingGateway(httpClient)
	imageRepository := repository.NewImageRepository(db)
	searchController := controller.NewSearchController(imageRepository, embeddingGateway)
	searchSimilarHandlerV1 := otelhttp.NewHandler(handlers.NewSearchSimilarHandler(searchController), "/")
	log.Println("bye")

	log.Println("About to create otelslog.NewLogger...")
	logger := otelslog.NewLogger("search")
	log.Println("xyz")
	
	slog.SetDefault(logger)
	log.SetOutput(os.Stderr) // eliminate circular log dependency error

	mux := http.NewServeMux()
	mux.Handle("/api/v1/searchSimilar", searchSimilarHandlerV1)


	server := &http.Server{
		Addr: ":8080",
		Handler: mux,
	}
	shutdownHooks = append(shutdownHooks, server.Shutdown)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		logger.Info("Starting server on port 8080")
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error(fmt.Sprintf("HTTP Server error: %v", err))
		}
		logger.Info("Stopped server")
	}()

	<-stop // blocks until stop os.Interrupt or syscall.SIGTERM is sent to process
	logger.Info("Attempting server shutdown...")
	
	if err := handleShutdownHooks(shutdownHooks); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	logger.Info("Goodbye")
}

func handleShutdownHooks(shutdownHooks []func(context.Context) error) error {
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	var errs error
	for _, shutdownHook := range shutdownHooks {
		if err := shutdownHook(shutdownCtx); err != nil {
			errors.Join(errs, err)
		}
	}
	return errs
}

