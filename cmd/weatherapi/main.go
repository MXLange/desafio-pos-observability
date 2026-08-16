package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MXLange/desafio-pos-observability/internal/clients/viacep"
	weatherclient "github.com/MXLange/desafio-pos-observability/internal/clients/weatherapi"
	"github.com/MXLange/desafio-pos-observability/internal/telemetry"
	weatherapiapp "github.com/MXLange/desafio-pos-observability/internal/weatherapi"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func main() {
	ctx := context.Background()

	shutdownTracer, err := telemetry.InitTracer(ctx, "weatherapi")
	if err != nil {
		log.Fatalf("failed to initialize tracer: %v", err)
	}

	apiKey := os.Getenv("WEATHER_API_KEY")
	if apiKey == "" {
		log.Fatal("WEATHER_API_KEY environment variable is required")
	}

	tracer := telemetry.Tracer("weatherapi")
	httpClient := &http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
		Timeout:   10 * time.Second,
	}

	service := weatherapiapp.NewService(
		viacep.NewClient(httpClient, tracer),
		weatherclient.NewClient(httpClient, apiKey, tracer),
	)
	handler := weatherapiapp.NewHandler(service)

	mux := http.NewServeMux()
	mux.Handle("/weather", otelhttp.NewHandler(handler, "weatherapi.handle-weather"))
	mux.Handle("/", otelhttp.NewHandler(handler, "weatherapi.handle-root"))

	port := envOrDefault("PORT", "8081")
	server := &http.Server{
		Addr:              ":" + port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("weatherapi listening on :%s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("server shutdown error: %v", err)
	}

	if err := shutdownTracer(ctx); err != nil {
		log.Printf("tracer shutdown error: %v", err)
	}
}

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
