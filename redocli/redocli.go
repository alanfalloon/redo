package main

import (
	"encoding/json"
	"fmt"
	"github.com/alanfalloon/redo/util"
	"log"
	"os"
)

func main() {
	log.SetPrefix(fmt.Sprint("redocli(", os.Getpid(), "): "))

	conn, err := util.Connect()
	if err != nil {
		_, conn = util.Launch("redod", nil, "", os.Stdout)
	}

	conn.Write(marshal(request()))
	display(response(conn))
}

func marshal(req interface{}) []byte {
	b, err := json.Marshal(req)
	if err != nil {
		log.Fatal("json: ", err)
	}
	return b
}

func response(conn *os.File) (resp util.Resp) {
	defer conn.Close()
	dec := json.NewDecoder(conn)
	err := dec.Decode(&resp)
	if err != nil {
		log.Fatal("response decode error:", err)
	}
	return
}

func display(resp util.Resp) {
	out := log.New(os.Stderr, "", 0)
	for _, line := range resp.Errlines {
		out.Println(line)
	}
	os.Exit(resp.ExitCode)
}
