# Build stage
FROM golang:1.21.9-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY main.go .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o projectdgt .

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates && \
    addgroup -S appgroup && adduser -S -G appgroup -H -h /app appuser

WORKDIR /app

COPY --from=builder /app/projectdgt .

USER appuser

HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 CMD [ "/app/projectdgt", "--health-check" ]

ENTRYPOINT ["/app/projectdgt"]