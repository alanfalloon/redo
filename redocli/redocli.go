package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func main() {
	log.SetPrefix(fmt.Sprint("redocli(", os.Getpid(), "): "))

	conn := connect()

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
