package main

import (
	"context"
	"grpc-gateway-example/gen/go/todo/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewServer(db Database) *Server {
	return &Server{
		Database: db,
	}
}

type Server struct {
	todo.UnimplementedYourServiceServer
	Database
}

// GetAllTODO implements todo.YourServiceServer.
func (s *Server) GetAllTODO(ctx context.Context, req *todo.IDRequest) (*todo.RepeatedStringMessage, error) {
	val, err := s.Database.FindAll(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get todo: %v", err)
	}

	return &todo.RepeatedStringMessage{Value: val}, nil
}

func (s *Server) CreateTODO(ctx context.Context, val *todo.StringMessage) (*todo.StringMessage, error) {
	err := s.Database.Create(ctx, val.Value)

	return &todo.StringMessage{Value: val.Value}, err
}

func (s *Server) UpdateTODO(ctx context.Context, val *todo.StringMessage) (*todo.StringMessage, error) {
	if val.Id == nil {
		return nil, status.Errorf(codes.InvalidArgument, "missing ID in request")
	}

	err := s.Database.Update(ctx, Int32ToIntP(*val.Id), val.Value)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update todo: %v", err)
	}

	return &todo.StringMessage{Value: val.Value}, nil
}

func (s *Server) GetTODO(ctx context.Context, val *todo.IDRequest) (*todo.StringMessage, error) {
	data, err := s.Database.Find(ctx, int(val.Id))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "failed to get todo: %v", err)
	}

	return &todo.StringMessage{Value: data}, nil
}

func (s *Server) DeleteTODO(ctx context.Context, val *todo.IDRequest) (*todo.StringMessage, error) {
	err := s.Database.Delete(ctx, int(val.Id))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "failed to delete todo: %v", err)
	}

	return &todo.StringMessage{Value: "deleted"}, nil
}

var _ todo.YourServiceServer = (*Server)(nil)
