package util

import (
	"fmt"
	"os"
)

func Connect() (conn *os.File, err error) {
	var redo_fd uintptr
	var n int
	n, err = fmt.Sscan(os.Getenv("REDO_FD"), &redo_fd)
	if n == 0 {
		return
	}
	conn = os.NewFile(redo_fd, "REDO_FD")
	_, err = conn.Stat()
	return
}
