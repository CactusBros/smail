# Build stage
FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /smaila ./cmd/main.go

# Runtime stage
FROM alpine:3 AS stage

RUN apk update && apk add --no-cache ca-certificates

# Create non-root user
# RUN adduser -D -g '' appuser

WORKDIR /app

COPY --from=builder /smaila ./

# Optional: Copy static files, configs, etc. (but not .env)
# COPY static/ ./static/

# USER appuser

EXPOSE 8080

ENTRYPOINT ["./smaila"]
