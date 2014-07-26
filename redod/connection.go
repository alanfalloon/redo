package main

import (
	"fmt"
	"os"
)

func connection() (conn *os.File) {
	var redo_fd uintptr
	_, err := fmt.Sscan(os.Getenv("REDO_FD"), &redo_fd)
	check(err)
	newf := os.NewFile(redo_fd, "REDO_FD")
	_, err = newf.Stat()
	check(err)
	log.Print("connection: established")
	return newf
}

func check(err error) {
	if err != nil {
		log.Fatal("check:", err)
	}
}
