package main

import redod "github.com/alanfalloon/redo/redod"

import (
	"encoding/json"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"
)

func main() {
	req := redod.Req{env(), os.Args, cwd()}
	b, err := json.Marshal(req)
	if err != nil {
		log.Fatal(err)
	}

	conn := dialDaemon()
	conn.Write(b)
	conn.CloseWrite()
	defer conn.CloseRead()

	dec := json.NewDecoder(conn)
	var outlines []string
	err = dec.Decode(&outlines)
	if err != nil {
		log.Fatal(err)
	}

	out := log.New(os.Stdout, "", 0)
	for _, line := range outlines {
		out.Println(line)
	}
	err = dec.Decode(&outlines)
	if err != io.EOF {
		log.Fatal("Expected EOF not:", err)
	}
}

func dialDaemon() *net.UnixConn {
	var once sync.Once
	for {
		conn, err := net.Dial("unix", "foo")
		if operr, ok := err.(*net.OpError); ok {
			switch operr.Err {
			case syscall.ENOENT, syscall.ECONNREFUSED:
				once.Do(launchDaemon)
			default:
				log.Fatal(operr)
			}
		} else if err != nil {
			log.Fatal(err)
		} else {
			return conn.(*net.UnixConn)
		}
	}
}

func launchDaemon() {
	cmd := exec.Command("redod", os.Args...)
	if err := cmd.Run(); err != nil {
		log.Fatal("redod died:", err)
	}
}

func env() map[string]string {
	env := make(map[string]string)
	for _, e := range os.Environ() {
		kv := strings.SplitN(e, "=", 2)
		k, v := kv[0], kv[1]
		if len(k) > 5 && k[0:5] == "REDO_" {
			env[k] = v
		}
	}
	return env
}

func cwd() string {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return cwd
}
