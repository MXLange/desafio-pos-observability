package gateway

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/MXLange/desafio-pos-observability/internal/zipcode"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type Handler struct {
	client *Client
}

func NewHandler(client *Client) *Handler {
	return &Handler{client: client}
}

type requestPayload struct {
	CEP string `json:"cep"`
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var payload requestPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		zipcode.WriteError(w, http.StatusUnprocessableEntity, zipcode.ErrInvalidZipcode)
		return
	}

	cep, err := zipcode.Validate(payload.CEP)
	if err != nil {
		zipcode.WriteError(w, http.StatusUnprocessableEntity, zipcode.ErrInvalidZipcode)
		return
	}

	body, statusCode, err := h.client.Forward(r.Context(), cep)
	if err != nil {
		zipcode.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, _ = w.Write(body)
}

func WeatherAPIURL() string {
	if value := os.Getenv("WEATHER_API_URL"); value != "" {
		return value
	}
	return "http://localhost:8081/weather"
}

func NewInstrumentedHTTPClient() *http.Client {
	return &http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
		Timeout:   10 * time.Second,
	}
}
