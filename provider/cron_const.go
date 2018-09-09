package provider

/**
	Cron服务提供者结构
 */

type CronProvider struct {
}

/**
	任务类型常量
 */

const (
	START = iota + 1
	STOP
	GRACE
	RUN_ONCE
	Info
)

