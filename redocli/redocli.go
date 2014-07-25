package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	log.SetPrefix(fmt.Sprint("redocli(", os.Getpid(), "): "))
	req := request()
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
