package main

import (
	"context"
	"grpc-gateway-example/gen/go/todo/v1"
	"log"
	"net"
	"net/http"

	_ "grpc-gateway-example/doc/statik"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	database := NewDatabase()

	go runGatewayServer(database)
	runGRPCServer(database)
}

func runGRPCServer(db Database) {
	server := NewServer(db)

	grpcserver := grpc.NewServer()
	todo.RegisterYourServiceServer(grpcserver, server)
	reflection.Register(grpcserver)

	listener, err := net.Listen("tcp", "127.0.0.1:5001")
	if err != nil {
		log.Default().Fatalf("failed to listen: %v", err)
	}

	log.Default().Println("gRPC server started")
	err = grpcserver.Serve(listener)
	if err != nil {
		log.Default().Fatalf("failed to serve: %v", err)
	}
}

func runGatewayServer(db Database) {
	server := NewServer(db)

	grpcmux := runtime.NewServeMux()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := todo.RegisterYourServiceHandlerServer(ctx, grpcmux, server)
	if err != nil {
		log.Default().Fatalf("failed to register gateway server: %v", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcmux)

	fs := http.FileServer(http.Dir("./gen/openapiv2/"))
	mux.Handle("/swagger/", http.StripPrefix("/swagger/", fs))

	listener, err := net.Listen("tcp", "127.0.0.1:5002")
	if err != nil {
		log.Default().Fatalf("failed to listen: %v", err)
	}
	log.Default().Println("http server started")

	err = http.Serve(listener, mux)
	if err != nil {
		log.Default().Fatalf("failed to serve: %v", err)
	}
}
