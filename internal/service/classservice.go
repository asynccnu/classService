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
func (s *ClassServiceService) AddClassByHand(ctx context.Context, req *pb.AddClassByHandRequest) (*pb.AddClassByHandReply, error) {
	return &pb.AddClassByHandReply{}, nil
}
func (s *ClassServiceService) AddClassBySearch(ctx context.Context, req *pb.AddClassBySearchRequest) (*pb.AddClassBySearchReply, error) {
	return &pb.AddClassBySearchReply{}, nil
}
