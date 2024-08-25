package client

import (
	"classService/internal/biz"
	"classService/internal/logPrinter"
	"context"
	v1 "github.com/asynccnu/Muxi_ClassList/api/classer/v1"
	"github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/google/wire"
)

const CLASSLISTSERVICE = "discovery:///MuXi_ClassList"

var ProviderSet = wire.NewSet(NewClassListService, NewClient)

type ClassListService struct {
	cs  v1.ClasserClient
	log logPrinter.LogerPrinter
}

func (c ClassListService) GetAllSchoolClassInfos(ctx context.Context) ([]biz.ClassInfo, error) {
	resp, err := c.cs.GetAllClassInfo(ctx, &v1.GetAllClassInfoRequest{})
	if err != nil {
		c.log.FuncError(c.cs.GetAllClassInfo, err)
		return nil, err
	}
	var classInfos = make([]biz.ClassInfo, 0)
	for _, info := range resp.ClassInfos {
		classInfo := biz.ClassInfo{
			ID:           info.Id,
			Day:          info.Day,
			Teacher:      info.Teacher,
			Where:        info.Where,
			ClassWhen:    info.ClassWhen,
			WeekDuration: info.WeekDuration,
			Classname:    info.Classname,
			Credit:       info.Credit,
			Weeks:        info.Weeks,
			Semester:     info.Semester,
			Year:         info.Year,
		}
		classInfos = append(classInfos, classInfo)
	}
	return classInfos, nil
}

func (c ClassListService) AddClassInfoToClassListService(ctx context.Context, req *v1.AddClassRequest) (*v1.AddClassResponse, error) {
	resp, err := c.cs.AddClass(ctx, req)
	if err != nil {
		c.log.FuncError(c.cs.AddClass, err)
		return nil, err
	}
	return resp, nil

}

func NewClassListService(cs v1.ClasserClient, log logPrinter.LogerPrinter) *ClassListService {
	return &ClassListService{
		cs:  cs,
		log: log,
	}
}
func NewClient(r *etcd.Registry, logger log.Logger) (v1.ClasserClient, error) {
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint(CLASSLISTSERVICE), // 需要发现的服务，如果是k8s部署可以直接用服务器本地地址:9001，9001端口是需要调用的服务的端口
		grpc.WithDiscovery(r),
		grpc.WithMiddleware(
			tracing.Client(),
			recovery.Recovery(),
		),
	)
	if err != nil {
		log.NewHelper(logger).WithContext(context.Background()).Errorw("kind", "grpc-client", "reason", "GRPC_CLIENT_INIT_ERROR", "err", err)
		return nil, err
	}
	return v1.NewClasserClient(conn), nil
}
