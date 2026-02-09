# ===== Build stage =====
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/server
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o geofence-worker ./cmd/geofence-worker
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o migrate ./cmd/migrate
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o mqtt-publisher ./cmd/mqtt-publisher

# ===== Runtime stage =====
FROM alpine:3.19

WORKDIR /app

COPY --from=builder /app/server .
COPY --from=builder /app/geofence-worker .
COPY --from=builder /app/migrate .
COPY --from=builder /app/mqtt-publisher .
COPY migrations ./migrations

EXPOSE 8080

# default: API
CMD ["./server"]