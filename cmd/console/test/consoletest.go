package main

import (
	"context"
	"fmt"
	"gitee.com/vincent78/gcutil/cmd/console"
	bm "gitee.com/vincent78/gcutil/model"
	"gitee.com/vincent78/gcutil/utils/timeUtil"
	"time"
)

func main() {
	console, err := console.New(console.NewConsoleConfig())
	defer console.Stop(false)
	if err != nil {
		fmt.Println(fmt.Sprintf("create the console object error: %v", err))
	} else {
		console.RegisterHandler("ll", "connect", WSConnect)

		console.Welcome()
		console.Interactive()
		console.StopInteractive()
	}
}

func WSConnect(_ context.Context, ch chan string, params ...string) bm.RespModel {
	time.Sleep(3 * time.Second)
	ch <- timeUtil.NowStr()
	time.Sleep(1 * time.Second)
	return bm.SuccessResp()
}
