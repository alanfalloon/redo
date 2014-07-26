package main

import (
	"log"
	"os"
	"os/exec"
)

func launchDaemon(conn *os.File) *exec.Cmd {
	cmd := exec.Command("redod")
	cmd.Stderr = os.Stderr
	cmd.ExtraFiles = []*os.File{conn}
	if err := cmd.Start(); err != nil {
		log.Fatal("redod died:", err)
	}
	return cmd
}
