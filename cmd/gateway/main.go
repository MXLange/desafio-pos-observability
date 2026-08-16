package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MXLange/desafio-pos-observability/internal/gateway"
	"github.com/MXLange/desafio-pos-observability/internal/telemetry"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func main() {
	ctx := context.Background()

	shutdownTracer, err := telemetry.InitTracer(ctx, "gateway")
	if err != nil {
		log.Fatalf("failed to initialize tracer: %v", err)
	}

	client := gateway.NewClient(gateway.WeatherAPIURL(), gateway.NewInstrumentedHTTPClient())
	handler := gateway.NewHandler(client)

	mux := http.NewServeMux()
	mux.Handle("/weather", otelhttp.NewHandler(handler, "gateway.handle-weather"))
	mux.Handle("/", otelhttp.NewHandler(handler, "gateway.handle-root"))

	port := envOrDefault("PORT", "8080")
	server := &http.Server{
		Addr:              ":" + port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("gateway listening on :%s", port)
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
