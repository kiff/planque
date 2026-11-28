//go:build !unix

package browserworker

import "os/exec"

func ownGroup(*exec.Cmd) {}

func killGroup(cmd *exec.Cmd) {
	if cmd.Process != nil {
		_ = cmd.Process.Kill()
	}
}

// noCoreDumps: no core-file limit to set off unix.
func noCoreDumps() error { return nil }
