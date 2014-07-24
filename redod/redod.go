package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	_log "log"
	"net"
	"os"
	"sync"
)

var log *_log.Logger = _log.New(os.Stderr, fmt.Sprint("redod(", os.Getpid(), "): "), _log.LstdFlags)

func logWrap(prefix string, l *_log.Logger) *_log.Logger {
	return _log.New(os.Stderr, fmt.Sprint(l.Prefix(), prefix), l.Flags())
}

var quit chan<- bool

func main() {
	var q = make(chan bool)
	quit = q
	go listen()
	<-q
	os.Exit(0)
}

func listen() {
	log := logWrap("listen: ", log)
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
	log := logWrap("handle: ", log)

	connections.Add(1)
	defer connections.Done()
	once.Do(func() {
		go reaper()
	})

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

var connections sync.WaitGroup
var once sync.Once

func reaper() {
	log := logWrap("reaper: ", log)
	log.Print("Started")
	connections.Wait()
	log.Print("Done; exiting")
	quit <- true
}
