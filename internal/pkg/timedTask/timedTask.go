package timedTask

import (
	"context"
	clog "github.com/asynccnu/classService/internal/log"
	"github.com/google/wire"
	"github.com/robfig/cron/v3"
)

var ProviderSet = wire.NewSet(NewTask)

// OptClassInfoToEs 定义接口
type OptClassInfoToEs interface {
	AddClassInfosToES(ctx context.Context, xnm, xqm string)
	DeleteSchoolClassInfosFromES(ctx context.Context, xnm, xqm string)
}

// Task 定义 Task 结构体
type Task struct {
	a OptClassInfoToEs
}

func NewTask(a OptClassInfoToEs) *Task {
	return &Task{
		a: a,
	}
}

// AddClassInfosToES 实现 Task 的 AddClassInfosToES 方法
func (t Task) AddClassInfosToES() {
	ctx := context.Background()

	StartAShortTask(func() {
		clog.LogPrinter.Info("开始执行 AddClassInfosToES 任务")
		t.a.AddClassInfosToES(ctx)
	})
}
func (t Task) DeleteSchoolClassInfosFromES() {
	ctx := context.Background()

	StartLongTimeTask(func() {
		clog.LogPrinter.Info("开始执行 DeleteSchoolClassInfosFromES 任务")
		t.a.DeleteSchoolClassInfosFromES(ctx)
	})
}

// StartAShortTask 用于启动定时任务
func StartAShortTask(task func()) {
	// 创建 Cron 实例
	c := cron.New()

	// 添加定时任务：每天凌晨 3 点执行
	c.AddFunc("0 0 3 * * *", task)
	//task()
	// 启动定时任务调度器
	c.Start()
}
func StartLongTimeTask(task func()) {
	// 创建 Cron 实例
	c := cron.New()

	// 添加定时任务：每隔3个月的1号凌晨3点执行
	c.AddFunc("0 0 3 */3 * *", task)

	// 启动定时任务调度器
	c.Start()
}
