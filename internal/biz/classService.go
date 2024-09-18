package biz

import (
	"context"
	v1 "github.com/asynccnu/be-api/gen/proto/classlist/v1"
	"github.com/asynccnu/classService/internal/logPrinter"
	"github.com/asynccnu/classService/internal/pkg/tool"
)

type EsProxy interface {
	AddClassInfo(ctx context.Context, classInfo ClassInfo) error
	RemoveClassInfo(ctx context.Context, xnm, xqm string)
	SearchClassInfo(ctx context.Context, keyWords string, xnm, xqm string) ([]ClassInfo, error)
}

type ClassListService interface {
	GetAllSchoolClassInfos(ctx context.Context, xnm, xqm string) ([]ClassInfo, error)
	AddClassInfoToClassListService(ctx context.Context, req *v1.AddClassRequest) (*v1.AddClassResponse, error)
}
type ClassSerivceUserCase struct {
	es  EsProxy
	cs  ClassListService
	log logPrinter.LogerPrinter
}

func NewClassSerivceUserCase(es EsProxy, cs ClassListService, printer logPrinter.LogerPrinter) *ClassSerivceUserCase {
	return &ClassSerivceUserCase{
		es:  es,
		cs:  cs,
		log: printer,
	}
}

func (c *ClassSerivceUserCase) AddClassInfoToClassListService(ctx context.Context, request *v1.AddClassRequest) (*v1.AddClassResponse, error) {
	return c.cs.AddClassInfoToClassListService(ctx, request)
}

func (c *ClassSerivceUserCase) SearchClassInfo(ctx context.Context, keyWords string, xnm, xqm string) ([]ClassInfo, error) {
	return c.es.SearchClassInfo(ctx, keyWords, xnm, xqm)
}

func (c *ClassSerivceUserCase) AddClassInfosToES(ctx context.Context) {
	xnm, xqm := tool.GetXnmAndXqm()
	classInfos, err := c.cs.GetAllSchoolClassInfos(ctx, xnm, xqm)
	if err != nil {
		c.log.FuncError(c.cs.GetAllSchoolClassInfos, err)
	}
	for _, classInfo := range classInfos {
		err1 := c.es.AddClassInfo(ctx, classInfo)
		if err1 != nil {
			c.log.FuncError(c.es.AddClassInfo, err1)
		}
	}
}
func (c *ClassSerivceUserCase) DeleteSchoolClassInfosFromES(ctx context.Context) {
	xnm, xqm := tool.GetXnmAndXqm()
	c.es.RemoveClassInfo(ctx, xnm, xqm)
}
