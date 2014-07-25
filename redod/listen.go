package main

import (
	"net"
	"os"
)

var listener net.Listener

func listen() {
	log := logWrap("listen: ", log)
	var err error
	listener, err = net.Listen("unix", "foo")
	if err != nil {
		log.Fatal(err)
	}

	// This signals to the redocli that launched us that we have
	// bound to the socket
	os.Stdout.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Got connection")
		go handle(conn)
	}
}
