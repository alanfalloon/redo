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
	go handle(conn, "")
	<-q
}

func check(err error) {
	util.Check(err)
}
