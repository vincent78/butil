//go:build darwin
// +build darwin

package exec

import (
	"context"
	"os"
	"os/exec"
)

func run(name string) *exec.Cmd {
	if len(EnvPath) > 0 {
		os.Setenv("PATH", EnvPath)
	}
	return exec.Command("/bin/sh", "-c", name)
}

func runWithCtx(ctx context.Context, cmd string) *exec.Cmd {
	return exec.CommandContext(ctx, "/bin/sh", "-c", cmd)
}
