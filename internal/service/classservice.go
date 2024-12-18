package service

import (
	"context"
	v1 "github.com/asynccnu/be-api/gen/proto/classlist/classlist"
	pb "github.com/asynccnu/classService/api/classService/v1"
	"github.com/asynccnu/classService/internal/biz"
	"github.com/asynccnu/classService/internal/logPrinter"
)

type ClassInfoProxy interface {
	AddClassInfoToClassListService(ctx context.Context, request *v1.AddClassRequest) (*v1.AddClassResponse, error)
	SearchClassInfo(ctx context.Context, keyWords string, xnm, xqm string) ([]biz.ClassInfo, error)
}

type ClassServiceService struct {
	pb.UnimplementedClassServiceServer
	cp  ClassInfoProxy
	log logPrinter.LogerPrinter
}

func NewClassServiceService(cp ClassInfoProxy, log logPrinter.LogerPrinter) *ClassServiceService {
	return &ClassServiceService{
		cp:  cp,
		log: log,
	}
}

func (s *ClassServiceService) SearchClass(ctx context.Context, req *pb.SearchRequest) (*pb.SearchReply, error) {
	classInfos, err := s.cp.SearchClassInfo(ctx, req.GetSearchKeyWords(), req.GetYear(), req.GetSemester())
	if err != nil {
		s.log.FuncError(s.cp.SearchClassInfo, err)
		return &pb.SearchReply{}, err
	}
	var pClassInfos = make([]*pb.ClassInfo, 0)
	for _, classInfo := range classInfos {
		info := HandleClassInfo(classInfo)
		pClassInfos = append(pClassInfos, info)
	}
	return &pb.SearchReply{
		ClassInfos: pClassInfos,
	}, nil
}

func (s *ClassServiceService) AddClass(ctx context.Context, req *pb.AddClassRequest) (*pb.AddClassReply, error) {
	preq := &v1.AddClassRequest{
		StuId:    req.GetStuId(),
		Name:     req.GetName(),
		DurClass: req.GetDurClass(),
		Where:    req.GetWhere(),
		Teacher:  req.GetTeacher(),
		Weeks:    req.GetWeeks(),
		Semester: req.GetSemester(),
		Year:     req.GetYear(),
		Day:      req.GetDay(),
	}
	if req.Credit != nil {
		*preq.Credit = req.GetCredit()
	}
	resp, err := s.cp.AddClassInfoToClassListService(ctx, preq)
	if err != nil {
		s.log.FuncError(s.cp.AddClassInfoToClassListService, err)
		return &pb.AddClassReply{}, err
	}
	return &pb.AddClassReply{
		Id:  resp.Id,
		Msg: resp.Msg,
	}, nil
}
func HandleClassInfo(classInfo biz.ClassInfo) *pb.ClassInfo {
	return &pb.ClassInfo{
		Day:          classInfo.Day,
		Teacher:      classInfo.Teacher,
		Where:        classInfo.Where,
		ClassWhen:    classInfo.ClassWhen,
		WeekDuration: classInfo.WeekDuration,
		Classname:    classInfo.Classname,
		Credit:       classInfo.Credit,
		Weeks:        classInfo.Weeks,
		Semester:     classInfo.Semester,
		Year:         classInfo.Year,
		Id:           classInfo.ID,
	}
}
