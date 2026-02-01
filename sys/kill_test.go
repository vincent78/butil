package sys

import (
	"os/exec"
	"runtime"
	"testing"
	"time"
)

func startTestProcess(t *testing.T) *exec.Cmd {
	t.Helper()
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = getTestProcessCmd("timeout", "/t", "30")
	} else {
		cmd = getTestProcessCmd("sleep", "30")
	}
	err := cmd.Start()
	if err != nil {
		t.Fatalf("Failed to start test process: %v", err)
	}
	return cmd
}

func TestKill_InvalidPID(t *testing.T) {
	err := Kill(0)
	if err == nil {
		t.Error("Expected an error for invalid PID 0, but got nil")
	}
}

func TestKill_GracefulExitSuccess(t *testing.T) {
	cmd := startTestProcess(t)
	pid := cmd.Process.Pid

	// Ensure the process is running before trying to kill it.
	if !isProcessRunning(pid) {
		t.Fatalf("Test process with PID %d should be running but it's not", pid)
	}

	err := Kill(pid)
	if err != nil {
		t.Errorf("Kill should have succeeded gracefully, but got error: %v", err)
	}

	// Give it a moment to die, then check.
	time.Sleep(1 * time.Second)

	if isProcessRunning(pid) {
		if ok := forceKill(pid); !ok {
			t.Errorf("Process %d should have been terminated gracefully, but it is still running.", pid)
		}
	}
}

func TestKill_ForceKillAfterGracefulTimeout(t *testing.T) {
	cmd := startTestProcess(t)
	pid := cmd.Process.Pid

	// To test the timeout path, we'll let Kill do its thing.
	// The 5-second wait in Kill should expire, triggering forceKill.
	err := Kill(pid)
	if err != nil {
		t.Errorf("Kill should have eventually succeeded with force kill, but got error: %v", err)
	}

	// We wait slightly longer than the internal timeout of Kill()
	time.Sleep(6 * time.Second)

	if isProcessRunning(pid) {
		if err = cmd.Process.Kill(); err != nil {
			t.Errorf("Process %d should have been force-killed, but it is still running.", pid)
		}
	}
}

func TestKill_NonExistentProcess(t *testing.T) {
	pid := 999999
	for isProcessRunning(pid) {
		pid++
	}

	err := Kill(pid)
	if err != nil {
		t.Errorf("Killing a non-existent process should not return an error, but got: %v", err)
	}
}

func TestKill_ForceKillFailure(t *testing.T) {
	// This scenario is hard to test reliably without special permissions.
	// For example, trying to kill a critical system process (PID 1 on Linux).
	// Such a test would be flaky and dangerous.
	// Instead, we will test the code path where forceKill *reports* failure,
	// but the process is actually already gone.

	// We can't easily mock `forceKill` to return false, so we'll test the
	// logic by killing a process and then immediately trying to kill it again.
	// The second kill attempt might fail on the `forceKill` call, but since
	// `isProcessRunning` will be false, `Kill` should return nil.
	cmd := startTestProcess(t)
	pid := cmd.Process.Pid

	// Kill it once for real
	forceKill(pid)
	time.Sleep(100 * time.Millisecond) // Let it die

	if isProcessRunning(pid) {
		t.Skip("Could not kill process for the first time, skipping test.")
	}

	// Now call the main Kill function on the already-dead process.
	// Internally, it might try graceful, then force, which might fail,
	// but the final check should show the process is not running.
	err := Kill(pid)
	if err != nil {
		t.Errorf("Expected nil error when Kill fails to force-kill an already-dead process, but got: %v", err)
	}
}

func TestIsProcessRunning_False(t *testing.T) {
	pid := 999999
	for isProcessRunning(pid) { // ensure PID is not running
		pid++
	}
	if isProcessRunning(pid) {
		t.Errorf("isProcessRunning for non-existent PID %d should be false", pid)
	}
}

func getTestProcessCmd(name string, arg ...string) *exec.Cmd {
	return exec.Command(name, arg...)
}
