FROM golang:1.26.5-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /weatherapi ./cmd/weatherapi

FROM alpine:3.22

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /weatherapi /app/weatherapi

EXPOSE 8081

ENTRYPOINT ["/app/weatherapi"]
