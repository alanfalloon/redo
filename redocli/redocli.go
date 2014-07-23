package main

import redod "github.com/alanfalloon/redo/redod"

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func main() {
	log.SetPrefix(fmt.Sprint("redocli(", os.Getpid(), "): "))
	req := redod.Req{env(), os.Args, cwd()}
	b, err := json.Marshal(req)
	if err != nil {
		log.Fatal("json: ", err)
	}

	conn, cmd := dialDaemon()
	defer waitDaemon(cmd)
	defer conn.CloseRead()
	conn.Write(b)
	conn.CloseWrite()

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
	os.Stdout.Close()
	err = dec.Decode(&outlines)
	if err != io.EOF {
		log.Fatal("Expected EOF not:", err)
	}
}

func dialDaemon() (conn *net.UnixConn, cmd *exec.Cmd) {
	for {
		_conn, err := net.Dial("unix", "foo")
		if operr, ok := err.(*net.OpError); ok {
			switch operr.Err {
			case syscall.ENOENT, syscall.ECONNREFUSED:
				if cmd == nil {
					cmd = launchDaemon()
				} else {
					log.Fatal("Already launched, but still: ", err)
				}
			default:
				log.Fatal("unexpected error: ", operr)
			}
		} else if err != nil {
			log.Fatal("unexpected non-operation error:", err)
		} else {
			conn = _conn.(*net.UnixConn)
			return
		}
	}
	panic("unreachable")
}

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
