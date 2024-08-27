package timedTask

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/robfig/cron/v3"
)

var ProviderSet = wire.NewSet(NewTask)

// AddClassInfoToEs 定义接口
type AddClassInfoToEs interface {
	AddClassInfosToES(ctx context.Context)
}

// Task 定义 Task 结构体
type Task struct {
	a AddClassInfoToEs
}

func NewTask(a AddClassInfoToEs) *Task {
	return &Task{
		a: a,
	}
}

// AddClassInfosToES 实现 Task 的 AddClassInfosToES 方法
func (t Task) AddClassInfosToES() {
	ctx := context.Background()
	log.Info("开始执行 AddClassInfosToES 任务")
	StartATask(func() {
		t.a.AddClassInfosToES(ctx)
	})
}

// StartATask 函数，用于启动定时任务
func StartATask(task func()) {
	// 创建 Cron 实例
	c := cron.New()

	// 添加定时任务：每天凌晨 3 点执行
	c.AddFunc("0 0 3 * * *", task)
	//task()
	// 启动定时任务调度器
	c.Start()
}
