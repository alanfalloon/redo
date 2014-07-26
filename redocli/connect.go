package main

import (
	"fmt"
	"log"
	"os"
	"syscall"
)

func connect() (conn *os.File) {
	var redo_fd uintptr
	n, err := fmt.Sscan(os.Getenv("REDO_FD"), &redo_fd)
	if n == 0 {
		return new_daemon()
	}
	newf := os.NewFile(redo_fd, "REDO_FD")
	_, err = newf.Stat()
	if err != nil {
		return new_daemon()
	}
	return newf
}

func new_daemon() (conn *os.File) {
	us, them := socketpair(syscall.AF_UNIX, syscall.SOCK_STREAM, 0)
	defer them.Close()
	err := os.Setenv("REDO_FD", "3")
	check(err)
	launchDaemon(them)
	return us
}

func check(err error) {
	if err != nil {
		log.Fatal("check:", err)
	}
}
