package viacep

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/MXLange/desafio-pos-observability/internal/zipcode"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

const defaultBaseURL = "https://viacep.com.br/ws"

type Client struct {
	baseURL    string
	httpClient *http.Client
	tracer     trace.Tracer
}

func NewClient(httpClient *http.Client, tracer trace.Tracer) *Client {
	return &Client{baseURL: defaultBaseURL, httpClient: httpClient, tracer: tracer}
}

type response struct {
	Localidade string          `json:"localidade"`
	Erro       json.RawMessage `json:"erro"`
}

func (r response) hasError() bool {
	if len(r.Erro) == 0 {
		return false
	}

	var asBool bool
	if err := json.Unmarshal(r.Erro, &asBool); err == nil {
		return asBool
	}

	var asString string
	if err := json.Unmarshal(r.Erro, &asString); err == nil {
		return asString == "true"
	}

	return false
}

func (c *Client) FetchCity(ctx context.Context, cep string) (string, error) {
	ctx, span := c.tracer.Start(ctx, "fetch-zipcode-location")
	defer span.End()

	span.SetAttributes(attribute.String("zipcode", cep))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/%s/json/", c.baseURL, cep), nil)
	if err != nil {
		recordError(span, err)
		return "", fmt.Errorf("create viacep request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		recordError(span, err)
		return "", fmt.Errorf("fetch zipcode location: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		recordError(span, err)
		return "", fmt.Errorf("read viacep response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("viacep returned status %d", resp.StatusCode)
		recordError(span, err)
		return "", err
	}

	var payload response
	if err := json.Unmarshal(body, &payload); err != nil {
		recordError(span, err)
		return "", fmt.Errorf("decode viacep response: %w", err)
	}

	if payload.hasError() || payload.Localidade == "" {
		recordError(span, zipcode.ErrNotFoundZipcode)
		return "", zipcode.ErrNotFoundZipcode
	}

	span.SetAttributes(attribute.String("city", payload.Localidade))
	return payload.Localidade, nil
}

func recordError(span trace.Span, err error) {
	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())
}
