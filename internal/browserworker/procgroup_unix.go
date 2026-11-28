//go:build unix

package browserworker

import (
	"os/exec"
	"syscall"
)

// ownGroup starts the browser in a process group of its own, so killGroup
// ends every process it started (renderers, the GPU and network services,
// the crash handler) and none is left writing to the profile.
func ownGroup(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
}

// killGroup kills the browser's whole process group.
func killGroup(cmd *exec.Cmd) {
	if cmd.Process != nil {
		_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
	}
}
