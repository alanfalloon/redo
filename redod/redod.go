package main

import (
	"github.com/alanfalloon/redo/util"
)

var quit chan<- bool

func main() {
	var q = make(chan bool)
	quit = q
	conn, err := util.Connect()
	util.Check(err)
	go handle(conn, "", -1)
	<-q
}

func check(err error, data ...interface{}) {
	if len(data) > 0 {
		log.Print(data...)
	}
	util.Check(err)
}
