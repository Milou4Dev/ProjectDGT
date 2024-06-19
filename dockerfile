FROM golang:1.21 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /app/main

FROM gcr.io/distroless/base-debian11

WORKDIR /app

COPY --from=builder /app/main .

USER nonroot:nonroot

CMD ["/app/main"]