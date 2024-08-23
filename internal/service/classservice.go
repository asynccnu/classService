package service

import (
	"context"

	pb "classService/api/classService/v1"
)

type ClassServiceService struct {
	pb.UnimplementedClassServiceServer
}

func NewClassServiceService() *ClassServiceService {
	return &ClassServiceService{}
}

func (s *ClassServiceService) SearchClass(ctx context.Context, req *pb.SearchRequest) (*pb.SearchReply, error) {
	return &pb.SearchReply{}, nil
}
func (s *ClassServiceService) AddClass(ctx context.Context, req *pb.AddClassRequest) (*pb.AddClassReply, error) {
	return &pb.AddClassReply{}, nil
}
