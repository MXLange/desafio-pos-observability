package weatherapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

const defaultBaseURL = "https://api.weatherapi.com/v1/current.json"

type Client struct {
	baseURL    string
	httpClient *http.Client
	apiKey     string
	tracer     trace.Tracer
}

func NewClient(httpClient *http.Client, apiKey string, tracer trace.Tracer) *Client {
	return &Client{baseURL: defaultBaseURL, httpClient: httpClient, apiKey: apiKey, tracer: tracer}
}

type response struct {
	Current struct {
		TempC float64 `json:"temp_c"`
	} `json:"current"`
}

func (c *Client) FetchTemperature(ctx context.Context, city string) (float64, error) {
	ctx, span := c.tracer.Start(ctx, "fetch-weather-temperature")
	defer span.End()

	span.SetAttributes(attribute.String("city", city))

	query := url.Values{}
	query.Set("key", c.apiKey)
	query.Set("q", city)
	query.Set("aqi", "no")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s?%s", c.baseURL, query.Encode()), nil)
	if err != nil {
		recordError(span, err)
		return 0, fmt.Errorf("create weatherapi request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		recordError(span, err)
		return 0, fmt.Errorf("fetch weather temperature: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		recordError(span, err)
		return 0, fmt.Errorf("read weatherapi response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("weatherapi returned status %d", resp.StatusCode)
		recordError(span, err)
		return 0, err
	}

	var payload response
	if err := json.Unmarshal(body, &payload); err != nil {
		recordError(span, err)
		return 0, fmt.Errorf("decode weatherapi response: %w", err)
	}

	span.SetAttributes(attribute.Float64("temp_c", payload.Current.TempC))
	return payload.Current.TempC, nil
}

func recordError(span trace.Span, err error) {
	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())
}
