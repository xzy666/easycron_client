package provider

import (
	"fmt"
	"runtime/debug"
	"time"
	"log"
)

type Job struct {
	Id          int                                               //Id
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
}

func (j *Job) Run() {
	if !j.Concurrent && j.Status > 0 {
		log.Fatal(fmt.Sprintf("任务[%d]上一次执行尚未结束，本次被忽略。", j.Id))
		return
	}

	defer func() {
		if err := recover(); err != nil {
			log.Fatal(err, "\n", string(debug.Stack()))
		}
	}()

	fmt.Sprintf("开始执行任务: %d", j.Id)

	j.Status++
	defer func() {
		j.Status--
	}()

	timeout := time.Duration(time.Hour * 24)
	if j.Timeout > 0 {
		timeout = time.Second * time.Duration(j.Timeout)
	}

	cmdOut, _, _, _ := j.RunFunc(timeout)
	fmt.Println(cmdOut)
}
