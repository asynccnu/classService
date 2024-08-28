//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"classService/internal/biz"
	"classService/internal/client"
	"classService/internal/conf"
	"classService/internal/data"
	"classService/internal/logPrinter"
	"classService/internal/pkg/timedTask"
	"classService/internal/registry"
	"classService/internal/server"
	"classService/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(*conf.Server, *conf.Data, *conf.Registry, log.Logger) (*APP, func(), error) {
	panic(wire.Build(server.ProviderSet,
		data.ProviderSet,
		biz.ProviderSet,
		service.ProviderSet,
		registry.ProviderSet,
		client.ProviderSet,
		timedTask.ProviderSet,
		logPrinter.ProviderSet,
		wire.Bind(new(biz.EsProxy), new(*data.Data)),
		wire.Bind(new(biz.ClassListSerivce), new(*client.ClassListService)),
		wire.Bind(new(timedTask.OptClassInfoToEs), new(*biz.ClassSerivceUserCase)),
		wire.Bind(new(service.ClassInfoProxy), new(*biz.ClassSerivceUserCase)),
		NewApp,
		newApp))
}
