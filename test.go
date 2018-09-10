package main

import (
	"net/rpc"
	"log"
	"fmt"
	cronProvider "easycron_client/provider/cron"
)

func main() {
	client, err := rpc.DialHTTP("tcp", "127.0.0.1:1234")
	if err != nil {
		log.Fatal("dialing:", err)
	}

	ct := cronProvider.CronTask{1, 4}
	var reply string
	err = client.Call("CronProvider.Cc", ct, &reply)
	if err != nil {
		log.Fatal("arith error:", err)
	}

	fmt.Println(reply)
}

