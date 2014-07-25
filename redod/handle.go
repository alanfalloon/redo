package main

import (
	"encoding/json"
	"io/ioutil"
	"net"
)

func handle(conn net.Conn) {
	defer conn.Close()
	defer forestall_reaping().Done()
	log := logWrap("handle: ", log)

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
	b, err = json.Marshal([]string{"ok"})
	if err != nil {
		log.Fatal(err)
	}
	conn.Write(b)
	log.Print("Done")
}
