FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY . ./
RUN go mod download
RUN go mod vendor

# Build the application
RUN go build -o main .

# Create the production image
FROM alpine
WORKDIR /app
COPY --from=builder /app/main .

EXPOSE 5002

# Replace with your actual health check script
HEALTHCHECK CMD ["/bin/sh", "-c", "echo -n 'OK'"]

CMD ["/app/main"]
