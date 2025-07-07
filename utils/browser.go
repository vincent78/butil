package utils

import (
	"os/exec"
	"runtime"
)

// AutoOpenBrowser auto open browser
func AutoOpenBrowser(visitURL string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}

	args = append(args, visitURL)
	return exec.Command(cmd, args...).Start()
}
