package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"sync"
)

var quit chan<- bool

func main() {
	log.SetPrefix(fmt.Sprint("redod(", os.Getpid(), "): "))
	var q = make(chan bool)
	quit = q
	go listen()
	<-q
	os.Exit(0)
}

func listen() {
	os.Remove("foo")
	listener, err := net.Listen("unix", "foo")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

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

type Req struct {
	Env  map[string]string "json:env"
	Argv []string          "json:argv"
	Cwd  string            "json:cwd"
}

func handle(conn net.Conn) {
	defer conn.Close()

	connections.Add(1)
	defer connections.Done()
	once.Do(func() {
		go reaper()
	})

	log.Print("reading")
	b, err := ioutil.ReadAll(conn)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("unmarshalling")
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
	log.Print("writing")
	conn.Write(b)
	log.Print("done")
}

var connections sync.WaitGroup
var once sync.Once

func reaper() {
	log.Print("Reaper started")
	connections.Wait()
	log.Print("Reaper done; exiting")
	quit <- true
}
