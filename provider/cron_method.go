package provider

import (
	"github.com/robfig/cron"
	"os/exec"
	"time"
	"fmt"
	"log"
)

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

