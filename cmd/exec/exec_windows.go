//go:build windows
// +build windows

package exec

import (
	"context"
	"os/exec"
)

func run(name string) *exec.Cmd {
	return exec.Command("cmd", "/C", name)
}

func runWithCtx(ctx context.Context, cmd string) *exec.Cmd {
	return exec.CommandContext(ctx, "cmd", "/C", cmd)
}
