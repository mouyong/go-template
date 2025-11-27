package scheduler

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

var (
	// Scheduler 全局调度器实例
	Scheduler *cron.Cron
)

// Init 初始化调度器
func Init() error {
	// 创建调度器 (使用中国时区)
	location, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return fmt.Errorf("加载时区失败: %w", err)
	}

	Scheduler = cron.New(cron.WithLocation(location))
	fmt.Println("✅ Scheduler: 已初始化 (时区: Asia/Shanghai)")
	return nil
}

// Start 启动调度器
func Start() {
	if Scheduler != nil {
		Scheduler.Start()
		fmt.Println("✅ Scheduler: 已启动")
	}
}

// Stop 停止调度器
func Stop() {
	if Scheduler != nil {
		Scheduler.Stop()
		fmt.Println("⏹️  Scheduler: 已停止")
	}
}

// AddJob 添加定时任务
func AddJob(spec string, cmd func()) (cron.EntryID, error) {
	if Scheduler == nil {
		return 0, fmt.Errorf("调度器未初始化")
	}
	return Scheduler.AddFunc(spec, cmd)
}
