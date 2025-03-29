// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/asynccnu/classService/internal/biz"
	"github.com/asynccnu/classService/internal/client"
	"github.com/asynccnu/classService/internal/conf"
	"github.com/asynccnu/classService/internal/data"
	"github.com/asynccnu/classService/internal/pkg/timedTask"
	"github.com/asynccnu/classService/internal/registry"
	"github.com/asynccnu/classService/internal/server"
	"github.com/asynccnu/classService/internal/service"
	"github.com/go-kratos/kratos/v2/log"
)

import (
	_ "go.uber.org/automaxprocs"
)

// Injectors from wire.go:

// wireApp init kratos application.
func wireApp(confServer *conf.Server, confData *conf.Data, confRegistry *conf.Registry, logger log.Logger) (*APP, func(), error) {
	elasticClient, err := data.NewEsClient(confData)
	if err != nil {
		return nil, nil, err
	}
	classData, cleanup, err := data.NewClassData(elasticClient)
	if err != nil {
		return nil, nil, err
	}
	etcdRegistry := registry.NewRegistrarServer(confRegistry)
	classListService, err := client.NewClassListService(etcdRegistry)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	classSerivceUserCase := biz.NewClassSerivceUserCase(classData, classListService)
	classServiceService := service.NewClassServiceService(classSerivceUserCase)
	freeClassroomData := data.NewFreeClassroomData(elasticClient)
	cookieSvc, err := client.NewCookieSvc(etcdRegistry)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	freeClassroomBiz := biz.NewFreeClassroomBiz(classData, freeClassroomData, cookieSvc)
	freeClassroomSvc := service.NewFreeClassroomSvc(freeClassroomBiz)
	grpcServer := server.NewGRPCServer(confServer, classServiceService, freeClassroomSvc, logger)
	selectionUploader := service.NewSelectionUploader(freeClassroomBiz)
	httpServer := server.NewHTTPServer(confServer, selectionUploader)
	app := newApp(logger, grpcServer, httpServer, etcdRegistry)
	task := timedTask.NewTask(classSerivceUserCase, freeClassroomBiz)
	mainAPP := NewApp(app, task)
	return mainAPP, func() {
		cleanup()
	}, nil
}
