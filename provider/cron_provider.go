package provider

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
)

var mainCron *cron.Cron
var lock sync.Mutex
var db *gorm.DB
//初始化定时任务
func init() {
	db, _ = gorm.Open("mysql", "root:root@/easycron?charset=utf8&parseTime=True&loc=Local")
	//jobList := make([]*Job, 50)
	mainCron = cron.New()
	mainCron.Start()
}


func (c *CronProvider) Cc(ct CronTask, reply *string) error {
	task := models.Task{}
	db.Find(&task, ct.Id)
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
}
//执行一次性脚本
func once(job *Job, reply *string)  {
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
	},job.Title)
	fmt.Println("执行一次性脚本")
}
//停止定时任务
func stop(job *Job, reply *string) {
	entry := GetEntryById(job.Id)

	err := mainCron.Remove(entry.Name)
	if err != nil{
		log.Fatal(err)
	}
	fmt.Println("结束了《" + entry.Name + "》任务")
}
//其他信息处理
func back(job *Job, reply *string) {
	entry := GetEntryById(1)
	fmt.Println(entry.Job)
	fmt.Println(entry.Prev)
	fmt.Println(entry.Next)
	fmt.Println("back")
}

