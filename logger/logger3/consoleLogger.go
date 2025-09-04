package logger3

import (
	"context"
	"fmt"
	"runtime/debug"

	"github.com/vincent78/butil/cmd/console"
	"github.com/vincent78/butil/config"
	"github.com/vincent78/butil/thread"
)

type ConsoleLogger struct {
	LogModel
	ShowDate bool `json:"showDate"` // 是否显示日期
}

func NewConsoleLogger(ctx context.Context, config config.LoggerConfig) *ConsoleLogger {
	log := &ConsoleLogger{
		LogModel: LogModel{
			Level:     LoggerLeverl(config.Level + 1),
			MsgFormat: MsgFormatDefault,
			Queue:     make(chan *MsgModel, 1024),
			ctx:       ctx,
		},
		ShowDate: false,
	}
	err := log.InitWithConfig(config)
	if err != nil {
		return nil
	}
	thread.Go(func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Printf(console.Yellow("context cancel"))
				return
			case m := <-log.Queue:
				if m.Level == LevelError {
					fmt.Println(console.Red(m.Str()))
				} else if m.Level == LevelWarn {
					fmt.Println(console.Yellow(m.Str()))
				} else {
					fmt.Println(m.Str())
				}
				if m.Level == LevelError {
					fmt.Println(string(debug.Stack()))
				}
				m.SendBack()
			}
		}
	})
	return log
}

func (log *ConsoleLogger) InitWithConfig(config config.LoggerConfig) error {
	log.Level = LoggerLeverl(config.Level + 1)
	return nil
}
