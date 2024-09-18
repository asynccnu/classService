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
	"github.com/asynccnu/classService/internal/logPrinter"
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
	elasticClient, err := data.NewEsClient(confData, logger)
	if err != nil {
		return nil, nil, err
	}
	logerPrinter := logPrinter.NewLogger(logger)
	dataData, cleanup, err := data.NewData(confData, elasticClient, logerPrinter, logger)
	if err != nil {
		return nil, nil, err
	}
	etcdRegistry := registry.NewRegistrarServer(confRegistry, logger)
	classerClient, err := client.NewClient(etcdRegistry, logger)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	classListService := client.NewClassListService(classerClient, logerPrinter)
	classSerivceUserCase := biz.NewClassSerivceUserCase(dataData, classListService, logerPrinter)
	classServiceService := service.NewClassServiceService(classSerivceUserCase, logerPrinter)
	grpcServer := server.NewGRPCServer(confServer, classServiceService, logger)
	app := newApp(logger, grpcServer, etcdRegistry)
	task := timedTask.NewTask(classSerivceUserCase)
	mainAPP := NewApp(app, task)
	return mainAPP, func() {
		cleanup()
	}, nil
}
