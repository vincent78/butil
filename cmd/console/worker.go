package console

import (
	"context"
	"fmt"
	"os"
	"os/exec"
)

func WorkerDone(id int, str string, ctx *context.Context, ret *int) {
	cmd := exec.CommandContext(*ctx, "sh", "-c", str)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(), fmt.Sprintf("ENV_TEST=%d", id))

	if err := cmd.Run(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
	}
	if cmd.ProcessState.Exited() {
		*ret = cmd.ProcessState.ExitCode()
	}
}
