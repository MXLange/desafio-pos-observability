# Desafio Pos Go Expert - Observability

Sistema distribuido em Go para consultar clima por CEP com tracing distribuido via OpenTelemetry, OTEL Collector e Zipkin.

## Arquitetura

- Gateway: recebe `POST /weather`, valida o CEP e encaminha para o WeatherAPI.
- WeatherAPI: consulta ViaCEP, consulta WeatherAPI e retorna cidade e temperaturas em Celsius, Fahrenheit e Kelvin.
- OTEL Collector: recebe spans OTLP dos servicos.
- Zipkin: visualiza os traces.

## Requisitos

- Docker e Docker Compose
- Chave da WeatherAPI em `WEATHER_API_KEY`

## Executando

```bash
export WEATHER_API_KEY=sua_chave
docker compose up --build
```

## Requisicao

Envie a requisicao para o Gateway:

```bash
curl -i -X POST http://localhost:8080/weather \
  -H 'Content-Type: application/json' \
  -d '{"cep":"89874000"}'
```

Resposta de sucesso:

```json
{
  "city": "Maravilha",
  "temp_C": 28.5,
  "temp_F": 83.3,
  "temp_K": 301.65
}
```

Erros esperados:

- `422 {"message":"invalid zipcode"}` quando o CEP nao for string com exatamente 8 digitos.
- `404 {"message":"can not find zipcode"}` quando o CEP valido nao for encontrado.

## Zipkin

Acesse:

```text
http://localhost:9411
```

Depois de executar uma chamada, procure pelos servicos `gateway` e `weatherapi`. O trace deve mostrar o fluxo:

```text
request -> gateway -> weatherapi -> ViaCEP -> WeatherAPI
```

Spans manuais incluidos no WeatherAPI:

- `fetch-zipcode-location`
- `fetch-weather-temperature`
