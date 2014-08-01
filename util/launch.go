package util

import (
	"os"
	"os/exec"
	"syscall"
)

func Launch(exe string, args []string, cwd string) (cmd *exec.Cmd, conn *os.File) {
	cmd = exec.Command(exe, args...)
	cmd.Stderr = os.Stderr
	us, them := socketpair(syscall.AF_UNIX, syscall.SOCK_STREAM, 0)
	defer them.Close()
	err := os.Setenv("REDO_FD", "3")
	Check(err)
	cmd.ExtraFiles = []*os.File{them}
	err = cmd.Start()
	Check(err)
	return cmd, us
}
