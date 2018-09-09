package main

import (
	"net/rpc"
	"net"
	"log"
	"net/http"
	"easycron_client/provider"
)

//////////////////////
// 微服务RPC平台客户端 //
/////////////////////

func main() {
	//开启cron RPC服务

	rpc.Register(new(provider.CronProvider))
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":1234")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	http.Serve(l, nil)
}