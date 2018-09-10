package provider

import (
	"fmt"
	"runtime/debug"
	"time"
	"log"
	"github.com/robfig/cron"
	"os/exec"
)


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
/**
	根据Id获取相关任务
 */
func GetEntryById(id int) *cron.Entry {
	entries := mainCron.Entries()
	for _, e := range entries {
		if v, ok := e.Job.(*Job); ok {
			if v.Id == id {
				return e
			}
		}
	}
	return nil
}

/**
	获取所有任务
 */
func GetEntries(size int) []*cron.Entry {
	ret := mainCron.Entries()
	if len(ret) > size {
		return ret[:size]
	}
	return ret
}

func runCmdWithTimeout(cmd *exec.Cmd, timeout time.Duration) (error, bool) {
	done := make(chan error)
	go func() {
		done <- cmd.Wait()
	}()

	var err error
	select {
	case <-time.After(timeout):
		log.Fatal(fmt.Sprintf("任务执行时间超过%d秒，进程将被强制杀掉: %d", int(timeout/time.Second), cmd.Process.Pid))
		go func() {
			<-done // 读出上面的goroutine数据，避免阻塞导致无法退出
		}()
		if err = cmd.Process.Kill(); err != nil {
			log.Fatal(fmt.Sprintf("进程无法杀掉: %d, 错误信息: %s", cmd.Process.Pid, err))
		}
		return err, true
	case err = <-done:
		return err, false
	}
}


