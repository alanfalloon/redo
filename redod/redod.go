package main

var quit chan<- bool

func main() {
	var q = make(chan bool)
	quit = q
	go listen()
	<-q
	listener.Close()
}
