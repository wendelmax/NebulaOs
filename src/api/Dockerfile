# Build Stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy go.work and module files
COPY go.work go.work.sum ./
COPY src/api/go.mod src/api/go.sum ./src/api/
COPY src/providers/go.mod src/providers/go.sum ./src/providers/

# Download dependencies
RUN go mod download

# Copy source code
COPY src/api/ ./src/api/
COPY src/providers/ ./src/providers/

# Build the API
WORKDIR /app/src/api
RUN go build -o /nebula-api cmd/server/main.go

# Run Stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

COPY --from=builder /nebula-api /nebula-api

EXPOSE 8000

CMD ["/nebula-api"]
