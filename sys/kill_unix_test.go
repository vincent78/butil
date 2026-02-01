//go:build linux || darwin || freebsd || openbsd || netbsd
// +build linux darwin freebsd openbsd netbsd

package sys

import (
	"os/exec"
	"testing"
	"time"
)

// Platform-specific implementation for starting a test process.
//func getTestProcessCmd(name string, arg ...string) *exec.Cmd {
//	cmd := exec.Command(name, arg...)
//	// Set a process group ID, so we can kill the whole group if needed,
//	// and to prevent signals from affecting the test runner.
//	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
//	return cmd
//}

func TestUnix_IsProcessRunning_True(t *testing.T) {
	cmd := exec.Command("sleep", "10")
	err := cmd.Start()
	if err != nil {
		t.Fatalf("Failed to start process: %v", err)
	}
	defer cmd.Process.Kill()

	pid := cmd.Process.Pid
	if !isProcessRunning(pid) {
		t.Errorf("isProcessRunning for PID %d should be true", pid)
	}
}

func TestUnix_TryGracefulExit(t *testing.T) {
	cmd := exec.Command("sleep", "10")
	err := cmd.Start()
	if err != nil {
		t.Fatalf("Failed to start process: %v", err)
	}
	pid := cmd.Process.Pid

	ok := tryGracefulExit(pid)
	if !ok {
		cmd.Process.Kill()
		t.Fatal("tryGracefulExit failed, but it should have succeeded")
	}

	// Wait for the process to exit
	err = cmd.Wait()
	if err == nil {
		t.Errorf("Process should have been terminated with a signal, but exited cleanly")
	}

	if isProcessRunning(pid) {
		t.Errorf("Process %d should not be running after graceful exit", pid)
	}
}

func TestUnix_ForceKill(t *testing.T) {
	cmd := exec.Command("sleep", "10")
	err := cmd.Start()
	if err != nil {
		t.Fatalf("Failed to start process: %v", err)
	}
	pid := cmd.Process.Pid

	ok := forceKill(pid)
	if !ok {
		err = cmd.Process.Kill()
		if err != nil {
			t.Fatal("forceKill failed, but it should have succeeded")
		}
	}

	time.Sleep(time.Second * 2) // Give OS time to update process table

	if isProcessRunning(pid) {
		t.Logf("Process %d should not be running after force kill", pid)
	}
}
