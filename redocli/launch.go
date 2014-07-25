package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

func launchDaemon() *exec.Cmd {
	cmd := exec.Command("redod", os.Args...)
	r, w, err := os.Pipe()
	if err != nil {
		log.Fatal(err)
	}
	cmd.Stdout = w
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		log.Fatal("redod died:", err)
	}
	w.Close()
	_, err = ioutil.ReadAll(r)
	if err != nil {
		log.Fatal(err)
	}
	return cmd
}

func waitDaemon(cmd *exec.Cmd) {
	if cmd == nil {
		return
	}
	err := cmd.Wait()
	if err != nil {
		log.Fatal(err)
	}
}
