package main

import (
	"encoding/json"
	"io"
	"log"
	"net"
)

func main() {
	listener, err := net.Listen("unix", "foo")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Got connection")
		go handle(conn)
	}
}

type Req struct {
	Cmd   string   "json:cmd"
	Cwd   string   "json:cwd"
	Files []string "json:files"
}

func handle(conn net.Conn) {
	defer conn.Close()
	dec := json.NewDecoder(conn)
	var req Req
	for {
		if err := dec.Decode(&req); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		log.Println("Recieved", req)
		b, err := json.Marshal(true)
		if err != nil {
			log.Fatal(err)
		}
		conn.Write(b)
	}
}
