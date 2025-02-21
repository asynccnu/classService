package biz

import (
	"context"
	"fmt"
	"github.com/asynccnu/classService/internal/client"
	"github.com/asynccnu/classService/internal/conf"
	"github.com/asynccnu/classService/internal/data"
	"github.com/asynccnu/classService/internal/registry"
	"testing"
)

var cs *ClassSerivceUserCase

func TestMain(m *testing.M) {
	cli, err := data.NewEsClient(&conf.Data{Es: &conf.Data_ES{
		Url:      "http://127.0.0.1:9200",
		Setsniff: false,
		Username: "elastic",
		Password: "12345678",
	}})
	if err != nil {
		panic(fmt.Sprintf("failed to create elasticsearch client: %v", err))
	}
	dt, _, _ := data.NewData(cli)
	etcdRegistry := registry.NewRegistrarServer(&conf.Registry{
		Etcd: &conf.Etcd{
			Addr:     "127.0.0.1:2379",
			Username: "",
			Password: "",
		},
	})
	classerClient, err := client.NewClient(etcdRegistry)
	if err != nil {
		return
	}
	classListService := client.NewClassListService(classerClient)

	cs = NewClassSerivceUserCase(dt, classListService)
	m.Run()
}

func TestClassSerivceUserCase_AddClassInfosToES(t *testing.T) {
	cs.AddClassInfosToES(context.Background(), "2024", "1")
}
