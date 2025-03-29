//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

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
		wire.Bind(new(biz.EsProxy), new(*data.ClassData)),
		wire.Bind(new(biz.ClassListService), new(*client.ClassListService)),
		wire.Bind(new(biz.FreeClassRoomData), new(*data.FreeClassroomData)),
		wire.Bind(new(biz.ClassData), new(*data.ClassData)),
		wire.Bind(new(biz.CookieClient), new(*client.CookieSvc)),
		wire.Bind(new(timedTask.ClassroomTask), new(*biz.FreeClassroomBiz)),
		wire.Bind(new(timedTask.OptClassInfoToEs), new(*biz.ClassSerivceUserCase)),
		wire.Bind(new(service.ClassInfoProxy), new(*biz.ClassSerivceUserCase)),
		wire.Bind(new(service.FreeClassRoomSaver), new(*biz.FreeClassroomBiz)),
		wire.Bind(new(service.FreeClassroomSearcher), new(*biz.FreeClassroomBiz)),
		NewApp,
		newApp))
}
