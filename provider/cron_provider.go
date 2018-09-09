package provider

import (
	"github.com/robfig/cron"
	"fmt"
	"time"
	"bytes"
	"os/exec"
	"log"
	"sync"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"easycron_client/models"
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

type CronTask struct {
	Id   int
	Type int
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
	case GRACE:
		go grace(job, reply)
	case RUN_ONCE:
		go once(job, reply)
	default:
		go back(job, reply)
	}
	///3.result ...
	return nil
}

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

func stop(job *Job, reply *string) {
	entry := GetEntryById(job.Id)

	err := mainCron.Remove(entry.Name)
	if err != nil{
		log.Fatal(err)
	}
	fmt.Println("结束了《" + entry.Name + "》任务")
}

func grace(job *Job, reply *string) {
	fmt.Println("重启任务")
}

func back(job *Job, reply *string) {
	entry := GetEntryById(1)
	fmt.Println(entry.Job)
	fmt.Println(entry.Prev)
	fmt.Println(entry.Next)
	fmt.Println("back")
}

func once(job *Job, reply *string) {
}
func AddJob(spec string, job *Job) bool {
	lock.Lock()
	defer lock.Unlock()

	if GetEntryById(job.Id) != nil {
		return false
	}
	err := mainCron.AddJob(spec, job, job.Title)
	if err != nil {
		log.Fatal("AddJob: ", err.Error())
		return false
	}
	return true
}
