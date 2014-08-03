package util

import (
	"io"
	"os"
	"os/exec"
	"syscall"
)

func Launch(exe string, args []string, cwd string, out io.Writer) (cmd *exec.Cmd, conn *os.File) {
	cmd = exec.Command(exe, args...)
	cmd.Stdout = out
	cmd.Stderr = os.Stderr
	cmd.Dir = cwd
	us, them := socketpair(syscall.AF_UNIX, syscall.SOCK_STREAM, 0)
	defer them.Close()
	err := os.Setenv("REDO_FD", "3")
	Check(err)
	cmd.ExtraFiles = []*os.File{them}
	err = cmd.Start()
	Check(err)
	return cmd, us
}
