//go:build freebsd
// +build freebsd

package exec

import (
	"os/exec"
)

func run(name string) *Cmd {
	return exec.Command("/bin/sh", "-c", name)
}

func runWithCtx(ctx context.Context, cmd string) *exec.Cmd {
	return exec.CommandContext(ctx, "/bin/sh", "-c", cmd)
}
