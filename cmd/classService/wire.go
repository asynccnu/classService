//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"classService/internal/biz"
	"classService/internal/conf"
	"classService/internal/data"
	"classService/internal/pkg/timedTask"
	"classService/internal/registry"
	"classService/internal/server"
	"classService/internal/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(*conf.Server, *conf.Data, *conf.Registry, log.Logger) (*kratos.App, *timedTask.Task, func(), error) {
	panic(wire.Build(server.ProviderSet,
		data.ProviderSet,
		biz.ProviderSet,
		service.ProviderSet,
		registry.ProviderSet,
		timedTask.ProviderSet,
		wire.Bind(new(biz.EsProxy), new(*data.Data)),
		wire.Bind(new(timedTask.AddClassInfoToEs), new(*biz.ClassSerivceUserCase)),
		newApp))
}
