package timedTask

import (
	"context"
	clog "github.com/asynccnu/classService/internal/log"
	"github.com/asynccnu/classService/internal/pkg/tool"
	"github.com/google/wire"
	"github.com/robfig/cron/v3"
	"time"
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
	c *cron.Cron
}

func NewTask(a OptClassInfoToEs) *Task {
	return &Task{
		a: a,
		c: cron.New(),
	}
}

// AddClassInfosToES 实现 Task 的 AddClassInfosToES 方法
func (t Task) AddClassInfosToES() {
	ctx := context.Background()

	err := t.startAShortTask(func() {
		clog.LogPrinter.Info("开始执行 AddClassInfosToES 任务")
		xnm, xqm := tool.GetXnmAndXqm(time.Now())
		t.a.AddClassInfosToES(ctx, xnm, xqm)
	})
	if err != nil {
		panic(err)
	}
}
func (t Task) DeleteSchoolClassInfosFromES() {
	ctx := context.Background()

	err := t.startLongTimeTask(func() {
		clog.LogPrinter.Info("开始执行 DeleteSchoolClassInfosFromES 任务")
		xnm, xqm := tool.GetXnmAndXqm(time.Now())
		t.a.DeleteSchoolClassInfosFromES(ctx, xnm, xqm)
	})
	if err != nil {
		panic(err)
	}
}

// startAShortTask 用于启动定时任务
func (t Task) startAShortTask(task func()) error {

	// 添加定时任务：每天凌晨 3 点执行
	_, err := t.c.AddFunc("0 3 * * *", task) // 每天凌晨 3 点执行

	if err != nil {
		clog.LogPrinter.Errorf("failed to add  short task")
		return err
	}
	//task()
	// 启动定时任务调度器
	t.c.Start()
	return nil
}
func (t Task) startLongTimeTask(task func()) error {

	// 每隔3个月的1号凌晨3点执行（5字段格式）
	_, err := t.c.AddFunc("0 3 1 */3 *", task)

	if err != nil {
		clog.LogPrinter.Errorf("failed to add long task")
		return err
	}

	// 启动定时任务调度器
	t.c.Start()
	return nil
}
