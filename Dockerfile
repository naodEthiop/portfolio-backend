FROM golang:1.24-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o portfolio-api ./cmd/api

FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /app

COPY --from=builder /app/portfolio-api /app/portfolio-api
COPY --from=builder /app/configs /app/configs
COPY --from=builder /app/uploads /app/uploads

EXPOSE 8080
ENTRYPOINT ["/app/portfolio-api"]
