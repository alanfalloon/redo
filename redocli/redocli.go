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
		_, conn = util.Launch("redod", nil, "")
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

func response(conn *os.File) []string {
	dec := json.NewDecoder(conn)
	var outlines []string
	err := dec.Decode(&outlines)
	if err != nil {
		log.Fatal("response decode error:", err)
	}
	return outlines
}

func display(lines []string) {
	out := log.New(os.Stdout, "", 0)
	for _, line := range lines {
		out.Println(line)
	}
}
