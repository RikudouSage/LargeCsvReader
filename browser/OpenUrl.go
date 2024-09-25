package browser

import (
	"os/exec"
	"runtime"
)

func OpenUrl(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
		break
	case "darwin":
		cmd = "open"
		break
	default:
		cmd = "xdg-open"
		break
	}

	args = append(args, url)

	return exec.Command(cmd, args...).Start()
}
