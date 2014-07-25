package main

import (
	"sync"
)

var connections sync.WaitGroup

func reaper() {
	log := logWrap("reaper: ", log)
	log.Print("Started")
	connections.Wait()
	log.Print("Done; exiting")
	quit <- true
}

var once sync.Once

func forestall_reaping() *sync.WaitGroup {
	connections.Add(1)
	once.Do(func() {
		go reaper()
	})
	return &connections
}
