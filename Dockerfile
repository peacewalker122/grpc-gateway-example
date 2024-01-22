FROM golang:1.21-alpine
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o main .
CMD ["/app/main"]

EXPOSE 5002
