// Package process provides functions to manage processes.
package sys

import (
	"fmt"
	"time"
)

// Kill terminates a process by its PID.
func Kill(pid int) error {
	if pid < 1 {
		return fmt.Errorf("invalid PID: %d", pid)
	}

	// 1. First, attempt a graceful shutdown.
	ok := tryGracefulExit(pid)
	if ok {
		// Wait for confirmation of exit: check every 0.25s for up to 10 times (5s total).
		waitInterval := 250 * time.Millisecond
		maxAttempts := 20
		for i := 0; i < maxAttempts; i++ {
			time.Sleep(waitInterval)
			if !isProcessRunning(pid) {
				return nil
			}
		}
	} else {
		// If the graceful signal failed, check if the process is already gone.
		if !isProcessRunning(pid) {
			return nil
		}
	}

	// 2. If graceful shutdown failed or timed out, force kill the process.
	ok = forceKill(pid)
	if !ok {
		// If force kill failed, do one last check to see if it's running.
		if !isProcessRunning(pid) {
			return nil
		}
		return fmt.Errorf("failed to terminate PID=%d (permissions or process does not exist)", pid)
	}
	return nil
}
