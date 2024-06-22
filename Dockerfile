FROM golang:1.21-alpine AS builder

WORKDIR /app

RUN apk add --no-cache ca-certificates gcc musl-dev tzdata

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w -extldflags '-static'" -o /app/main

RUN go test ./...

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

WORKDIR /app

COPY --from=builder /app/main .

USER 65534:65534

ENV TZ=UTC

EXPOSE 8080

ENTRYPOINT ["/app/main"]