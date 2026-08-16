package weatherapi

import (
	"context"

	"github.com/MXLange/desafio-pos-observability/internal/clients/viacep"
	weatherclient "github.com/MXLange/desafio-pos-observability/internal/clients/weatherapi"
	"github.com/MXLange/desafio-pos-observability/internal/weather"
	"github.com/MXLange/desafio-pos-observability/internal/zipcode"
)

type Service struct {
	viaCEP     *viacep.Client
	weatherAPI *weatherclient.Client
}

func NewService(viaCEP *viacep.Client, weatherAPI *weatherclient.Client) *Service {
	return &Service{viaCEP: viaCEP, weatherAPI: weatherAPI}
}

type Response struct {
	City  string  `json:"city"`
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

func (s *Service) GetWeatherByCEP(ctx context.Context, cep string) (*Response, error) {
	validCEP, err := zipcode.Validate(cep)
	if err != nil {
		return nil, err
	}

	city, err := s.viaCEP.FetchCity(ctx, validCEP)
	if err != nil {
		return nil, err
	}

	tempC, err := s.weatherAPI.FetchTemperature(ctx, city)
	if err != nil {
		return nil, err
	}

	tempF, tempK := weather.ConvertFromCelsius(tempC)

	return &Response{City: city, TempC: tempC, TempF: tempF, TempK: tempK}, nil
}
