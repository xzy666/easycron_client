package cron

import "time"

/**
	Cron服务提供者结构
 */

type CronProvider struct {
}

/**
	rpc客户端投递的任务结构
 */
type CronTask struct {
	Id   uint
	Type int
}

/**
	任务类型常量
 */

const (
	START    = iota + 1
	STOP
	RUN_ONCE
	INFO
)

type Job struct {
	Id          uint                                              //Id
	LogId       int                                               //日志Id
	Type        int                                               //任务类型 开启关闭重启CRON
	Title       string                                            //任务名称
	RunFunc     func(time.Duration) (string, string, error, bool) // 执行函数
	Status      int                                               // 任务状态，大于0表示正在执行中
	Concurrent  bool                                              // 同一个任务是否允许并行执行
	Description string                                            //任务描述
	Spec        string                                            //定时脚本时间表达式
	Command     string                                            //脚本
	Timeout     int                                               //超时时间
	Context     map[string]string                                 //脚本执行的上下文信息
	State       int                                               //任务目前是否在执行或预备当中
}
