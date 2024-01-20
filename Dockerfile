FROM golang:1.21.2-alpine
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o main .
CMD ["/app/main"]

