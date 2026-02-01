//go:build linux || darwin || freebsd || openbsd || netbsd
// +build linux darwin freebsd openbsd netbsd

package sys

import (
	"syscall"
)

// tryGracefulExit: sends SIGTERM, returns true if the signal was sent successfully (does not guarantee immediate exit).
// On Unix, sending SIGTERM is the conventional way to request a graceful shutdown.
func tryGracefulExit(pid int) bool {
	// A nil error from syscall.Kill indicates the signal was sent successfully
	// (i.e., the target process exists and received the signal).
	err := syscall.Kill(pid, syscall.SIGTERM)
	return err == nil
}

// forceKill: sends SIGKILL to forcibly terminate the process.
func forceKill(pid int) bool {
	err := syscall.Kill(pid, syscall.SIGKILL)
	return err == nil
}

// isProcessRunning: checks if a process exists using syscall.Kill(pid, 0).
func isProcessRunning(pid int) bool {
	err := syscall.Kill(pid, 0)
	// err == nil => process exists
	// err == ESRCH => process does not exist
	// err == EPERM => process exists, but we don't have permission to send a signal (we count this as running)
	if err == nil {
		return true
	}
	// If the error is "operation not permitted," it means the process exists, but we lack permissions.
	if err == syscall.EPERM {
		return true
	}
	// Other errors (like ESRCH) indicate the process does not exist or another issue occurred.
	return false
}
