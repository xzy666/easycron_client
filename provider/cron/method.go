package cron

import (
	"github.com/robfig/cron"
	"fmt"
	"log"
	"sync"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"easycron_client/models"
	"time"
	"bytes"
	"os/exec"
	"github.com/go-ini/ini"
	"os"
	"easycron_client/base"
	"encoding/json"
)

var mainCron *cron.Cron
var lock sync.Mutex
var db *gorm.DB

// 初始化定时任务
func init() {
	cfg, err := ini.Load("config/app.ini")
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}

	db, _ = gorm.Open("mysql", fmt.Sprintf(
		"%s:%s@/%s?charset=utf8&parseTime=True&loc=Local",
		cfg.Section("database").Key("username").String(),
		cfg.Section("database").Key("password").String(),
		cfg.Section("database").Key("database").String()))

	mainCron = cron.New()
	mainCron.Start()

	//如果重启了程序 则重新读取一下数据库任务 初始化任务
}

func (c *CronProvider) Cc(ct CronTask, reply *string) error {
	task := models.Task{}
	db.Find(&task, ct.Id)
	// 未找到相关的任务
	if task.ID == 0 {
		*reply = string(base.JsonMsg(500))
		return nil
	}
	///1.middle & init
	job := &Job{
		ct.Id,
		task.LogId,
		ct.Type,
		task.Title,
		nil,
		task.Status,
		task.Concurrent,
		task.Description,
		task.Spec,
		task.Command,
		task.Timeout,
		nil}
	///2.deliver Rpc
	switch job.Type {
	case START:
		go start(job, reply)
	case STOP:
		go stop(job, reply)
	case RUN_ONCE:
		go once(job, reply)
	case INFO:
		go info(job, reply)
	default:
		go back(job, reply)
	}
	///3.result ...
	return nil
}

//开启定时任务
func start(job *Job, reply *string) {
	job.RunFunc = func(duration time.Duration) (string, string, error, bool) {
		bufOut := new(bytes.Buffer)
		bufErr := new(bytes.Buffer)
		cmd := exec.Command("/bin/bash", "-c", job.Command)
		cmd.Stdout = bufOut
		cmd.Stderr = bufErr
		cmd.Start()
		err, isTimeout := runCmdWithTimeout(cmd, duration)

		return bufOut.String(), bufErr.String(), err, isTimeout
	}
	AddJob(job.Spec, job)
	log.Println("开启任务")
	*reply = string(base.JsonMsg(200))
}

//执行一次性脚本
func once(job *Job, reply *string) {
	job.RunFunc = func(duration time.Duration) (string, string, error, bool) {
		bufOut := new(bytes.Buffer)
		bufErr := new(bytes.Buffer)
		cmd := exec.Command("/bin/bash", "-c", job.Command)
		cmd.Stdout = bufOut
		cmd.Stderr = bufErr
		cmd.Start()
		err, isTimeout := runCmdWithTimeout(cmd, duration)
		fmt.Println(bufOut.String())
		return bufOut.String(), bufErr.String(), err, isTimeout
	}
	timeout := time.Duration(time.Hour * 24)
	if job.Timeout > 0 {
		timeout = time.Second * time.Duration(job.Timeout)
	}
	mainCron.AddOnceFunc(job.Spec, func() {
		job.RunFunc(timeout)
	}, job.Title)
	fmt.Println("执行一次性脚本")
	*reply = string(base.JsonMsg(200))
}

//停止定时任务
func stop(job *Job, reply *string) {
	entry := GetEntryById(job.Id)

	err := mainCron.Remove(entry.Name)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("结束了《" + entry.Name + "》任务")
}

//获取定时任务信息
func info(job *Job, reply *string) {
	entry := GetEntryById(job.Id)
	data, _ := json.Marshal(map[string]string{"nextTime": entry.Next.Format("2006-01-02 15:04:05"), "prevTime": entry.Prev.Format("2006-01-02 15:04:05"), "name": entry.Name})
	*reply = string(data)
	fmt.Println("获取定时任务的信息")
}


//其他信息处理
func back(job *Job, reply *string) {
	*reply = string(base.JsonMsg(400))
	fmt.Println("back")
}
