package thread

import (
	"fmt"
	"time"
)

var daemonQueue = make(map[string]Daemon)

// var logger = logs.GetLogger(logs.LG_DAEMON)
var tmp int

type Daemon interface {
	Done()
}

func init() {
	DaemonStart()
}

func DaemonStart() {
	go func() {
		//logger.Info("daomen is begined")
		t := time.NewTicker(time.Second * time.Duration(time.Second*5))
		go func() {
			for range t.C {
				doTask()
				if r := recover(); r != nil {
					fmt.Printf("daomen error: %v", r)
				}
			}
		}()
	}()
}

func doTask() {
	//logger.Info("goroutine:%v", runtime.NumGoroutine())
	//logger.Flush()
}
