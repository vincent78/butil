package exec

import (
	"context"
	"gitee.com/vincent78/gcutil/model"
	"testing"
	"time"
)

func TestExec(t *testing.T) {
	ShellCommand("ls -lh")
}

func TestExecWithPipe(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	ch := make(chan string)
	rt := make(chan model.RespModel)
	//go ShellCommandWithCtx(ctx, ch, rt, "pwd && ls -lh")
	go ShellCommandWithCtx(ctx, ch, rt, "/users/vincent/temp/test.sh")

	for {
		select {
		case str := <-ch:
			t.Log(str)
		case <-time.After(30 * time.Second):
			t.Error("the time out")
			cancel()
		case result := <-rt:
			if result.Code != model.Success {
				t.Error(result.Message)
			} else {
				t.Log("cmd process end")
			}
			return
		}
	}
}
