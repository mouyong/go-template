package scheduler

import (
	"fmt"
)

// Job 定时任务接口
type Job struct {
	Spec        string // Cron 表达式
	Handler     func() // 任务处理函数
	Description string // 任务描述
}

var registeredJobs []Job

// Register 注册定时任务 (在调度器启动前调用)
func Register(job Job) {
	registeredJobs = append(registeredJobs, job)
}

// LoadJobs 加载所有已注册的任务到调度器
func LoadJobs() error {
	if Scheduler == nil {
		return fmt.Errorf("调度器未初始化")
	}

	if len(registeredJobs) == 0 {
		fmt.Println("⚠️  没有注册任何定时任务")
		return nil
	}

	for _, job := range registeredJobs {
		if _, err := AddJob(job.Spec, job.Handler); err != nil {
			return fmt.Errorf("加载任务失败 [%s]: %w", job.Description, err)
		}
		fmt.Printf("✅ 定时任务: %s (%s)\n", job.Description, job.Spec)
	}

	return nil
}
