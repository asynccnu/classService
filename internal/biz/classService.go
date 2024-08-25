package biz

import (
	"classService/internal/logPrinter"
	"context"
	v1 "github.com/asynccnu/Muxi_ClassList/api/classer/v1"
)

type EsProxy interface {
	AddClassInfo(ctx context.Context, classInfo ClassInfo) error
	SearchClassInfo(ctx context.Context, keyWords string) ([]ClassInfo, error)
}
type ClassListSerivce interface {
	GetAllSchoolClassInfos(ctx context.Context) ([]ClassInfo, error)
	AddClassInfoToClassListService(ctx context.Context, req *v1.AddClassRequest) (*v1.AddClassResponse, error)
}
type ClassSerivceUserCase struct {
	es  EsProxy
	cs  ClassListSerivce
	log logPrinter.LogerPrinter
}

func NewClassSerivceUserCase(es EsProxy, cs ClassListSerivce, printer logPrinter.LogerPrinter) *ClassSerivceUserCase {
	return &ClassSerivceUserCase{
		es:  es,
		cs:  cs,
		log: printer,
	}
}
func (c *ClassSerivceUserCase) AddClassInfoToClassListService(ctx context.Context, request *v1.AddClassRequest) (*v1.AddClassResponse, error) {
	return c.cs.AddClassInfoToClassListService(ctx, request)
}
func (c *ClassSerivceUserCase) SearchClassInfo(ctx context.Context, keyWords string) ([]ClassInfo, error) {
	return c.es.SearchClassInfo(ctx, keyWords)
}
func (c *ClassSerivceUserCase) AddClassInfosToES(ctx context.Context) {
	classInfos, err := c.cs.GetAllSchoolClassInfos(ctx)
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
