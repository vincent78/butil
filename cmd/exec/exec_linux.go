//go:build linux
// +build linux

package exec

import (
	"context"
	"os/exec"
)

func run(name string) *exec.Cmd {
	return exec.Command("/bin/sh", "-c", name)
}

func runWithCtx(ctx context.Context, cmd string) *exec.Cmd {
	return exec.CommandContext(ctx, "/bin/sh", "-c", cmd)
}
