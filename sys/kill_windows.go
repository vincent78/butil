//go:build windows
// +build windows

package sys

import (
	"golang.org/x/sys/windows"
)

var (
	kernel32                  = windows.NewLazySystemDLL("kernel32.dll")
	procGenerateConsoleCtrl   = kernel32.NewProc("GenerateConsoleCtrlEvent")
	procAttachConsole         = kernel32.NewProc("AttachConsole")
	procFreeConsole           = kernel32.NewProc("FreeConsole")
	procSetConsoleCtrlHandler = kernel32.NewProc("SetConsoleCtrlHandler")
)

const (
	CTRL_C_EVENT     = 0
	CTRL_BREAK_EVENT = 1
	STILL_ACTIVE     = 259
	// ATTACH_PARENT_PROCESS defines the parameter for attaching to the parent process's console.
	// Its value is -1, which is ^uintptr(0) in Go for cross-architecture compatibility.
	ATTACH_PARENT_PROCESS = ^uintptr(0)
)

// tryGracefulExit: Attempts to attach to the target process's console and send a CTRL_BREAK_EVENT.
// This modified version ensures the calling process's state is restored after execution.
func tryGracefulExit(pid int) bool {
	// Ensure that when the function exits, for any reason, it re-attaches
	// to the parent process's console. This will restore the ability to continue
	//outputting logs in the terminal.
	defer procAttachConsole.Call(ATTACH_PARENT_PROCESS) //nolint

	// Detach from the current console. This is necessary to attach to the target process's console.
	// This is the direct cause of the temporary interruption of logging.
	_, _, _ = procFreeConsole.Call()

	// Attach to the target process's console.
	r, _, _ := procAttachConsole.Call(uintptr(pid))
	if r == 0 {
		// If attaching fails (e.g., the target process has no console), return immediately.
		// The deferred procAttachConsole.Call above will handle re-attaching to our own console.
		return false
	}
	// If attaching succeeds, it's crucial to detach from the target's console upon exiting the function.
	defer procFreeConsole.Call() //nolint

	// Temporarily ignore the signal
	_, _, _ = procSetConsoleCtrlHandler.Call(CTRL_C_EVENT, CTRL_BREAK_EVENT)
	// Use defer to ensure that default handling of console signals is restored before the function returns.
	// The parameter '0' (FALSE) removes the NULL handler, thus restoring the functionality of CTRL+C.
	defer procSetConsoleCtrlHandler.Call(CTRL_C_EVENT, CTRL_C_EVENT) //nolint

	// Send a CTRL_BREAK_EVENT signal to the process group of the target process.
	r, _, _ = procGenerateConsoleCtrl.Call(uintptr(CTRL_BREAK_EVENT), uintptr(pid))
	return r != 0
}

// forceKill: Forcibly terminates the process using TerminateProcess.
func forceKill(pid int) bool {
	h, err := windows.OpenProcess(windows.PROCESS_TERMINATE, false, uint32(pid))
	if err != nil {
		return false
	}
	defer windows.CloseHandle(h) //nolint
	err = windows.TerminateProcess(h, uint32(1))
	return err == nil
}

// isProcessRunning: Checks if the process is still running using OpenProcess and GetExitCodeProcess.
func isProcessRunning(pid int) bool {
	h, err := windows.OpenProcess(windows.PROCESS_QUERY_LIMITED_INFORMATION, false, uint32(pid))
	if err != nil {
		return false
	}
	defer windows.CloseHandle(h) //nolint
	var code uint32
	err = windows.GetExitCodeProcess(h, &code)
	if err != nil {
		return false
	}
	return code == STILL_ACTIVE
}
