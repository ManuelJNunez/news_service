# ) Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy project configuration files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy src code inside container
COPY . .

# Compile REST API application
RUN CGO_ENABLED=0 GOOS=linux go build -o api ./cmd/api

# 2) Production stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

# Create unprivileged user to run the application
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

WORKDIR /app

# Copy the built binary from the previous stage and HTML templates for HTTP responses
COPY --from=builder /app/api .
COPY templates/ templates/

# Make appuser the owner of the app directory
RUN chown -R appuser:appuser /app

# Activate the unprivileged user
USER appuser

EXPOSE 8000

CMD ["./api"]
