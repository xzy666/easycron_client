package main

import (
	"net/rpc"
	"net"
	"log"
	"net/http"
	cronProvider "easycron_client/provider/cron"
	"github.com/go-ini/ini"
	"fmt"
	"os"
)

//////////////////////
// 微服务RPC平台客户端 //
/////////////////////

func main() {
	//载入配置
	cfg, err := ini.Load("config/app.ini")
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}

	//开启cron RPC服务
	rpc.Register(new(cronProvider.CronProvider))
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":"+cfg.Section("rpc").Key("port").String())
	if e != nil {
		log.Fatal("listen error:", e)
	}
	http.Serve(l, nil)
}
