package console

import (
	"context"
	"strings"
	"sync"
	"testing"
)

func TestWorkerDone(t *testing.T) {
	var (
		wg  sync.WaitGroup
		ret int
	)

	args := []string{"echo  1", "echo  2", "echo  3"}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for wid, wargs := range args {
		wg.Add(1)
		go func(wid int, wargs string) {
			defer wg.Done()
			defer cancel()
			cmds := []string{"$ENV_TEST"}
			cmds = append(cmds)
			WorkerDone(wid, "echo $ENV_TEST", &ctx, &ret)
		}(wid, strings.TrimSpace(wargs))
	}

	wg.Wait()
	t.Logf("%d workers done", ret)
}
