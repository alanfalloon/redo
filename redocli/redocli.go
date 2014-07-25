package main

import redod "github.com/alanfalloon/redo/redod"

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
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
