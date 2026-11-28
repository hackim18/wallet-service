# Build stage
FROM golang:1.25-alpine AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/wallet-service ./cmd/web

# Runtime stage: a small image that runs as a non-root user
FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata && adduser -D -u 10001 app
WORKDIR /app
COPY --from=build /out/wallet-service /app/wallet-service
USER app

EXPOSE 8080
# Configuration comes from environment variables (see .env.example).
ENTRYPOINT ["/app/wallet-service"]
CMD ["--migrate", "--run"]
