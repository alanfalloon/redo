package main

import (
	"sync"
)

var connections sync.WaitGroup
var once sync.Once

func reaper() {
	log := logWrap("reaper: ", log)
	log.Print("Started")
	connections.Wait()
	log.Print("Done; exiting")
	quit <- true
}
