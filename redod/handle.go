package main

import (
	"encoding/json"
	"io/ioutil"
	"net"
)

func handle(conn net.Conn) {
	defer conn.Close()
	defer forestall_reaping().Done()

	log := logWrap("handle:", log)
	log.Print("begin", conn)
	defer log.Print("done", conn)

	var quit = make(chan bool)
	go reply(conn, fulfill(requests(conn)), quit)
	<-quit
}

type resp []string

func reply(conn net.Conn, resps <-chan resp, quit chan<- bool) {
	defer conn.Close()
	defer func() { quit <- true }()

	log := logWrap("reply:", log)
	log.Print("begin", conn)
	defer log.Print("done", conn)

	for resp := range resps {
		b, err := json.Marshal(resp)
		if err != nil {
			log.Fatal(err)
		}
		conn.Write(b)
	}
}

func requests(conn net.Conn) <-chan Req {
	var sink = make(chan Req, 1)
	go func(sink chan<- Req) {
		defer func() { close(sink) }()

		log := logWrap("requests:", log)
		log.Print("begin", conn)
		defer log.Print("done", conn)

		log.Print("Reading")
		b, err := ioutil.ReadAll(conn)
		if err != nil {
			log.Fatal(err)
		}
		var req Req
		err = json.Unmarshal(b, &req)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Recieved", req)
		sink <- req
	}(sink)
	return sink
}
