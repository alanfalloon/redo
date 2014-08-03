package main

import (
	"encoding/json"
	"io"
	"os"
)

func handle(conn *os.File, cwd string) {
	defer conn.Close()
	defer forestall_reaping().Done()
	reply(conn, fulfill(requests(conn), cwd))
}

type resp []string

func reply(conn *os.File, resps <-chan resp) {
	defer conn.Close()
	log := logWrap("reply:", log)
	for resp := range resps {
		b, err := json.Marshal(resp)
		if err != nil {
			log.Fatal(err)
		}
		conn.Write(b)
	}
}

func requests(conn *os.File) <-chan req {
	var sink = make(chan req, 1)
	go func(sink chan<- req) {
		defer func() { close(sink) }()
		log := logWrap("requests:", log)
		dec := json.NewDecoder(conn)
		var req req
		var err error
		for err = dec.Decode(&req); err == nil; err = dec.Decode(&req) {
			sink <- req
		}
		if err != io.EOF {
			log.Fatal("fatal recv error:", err)
		}
	}(sink)
	return sink
}
