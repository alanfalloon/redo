package main

var quit chan<- bool

func main() {
	var q = make(chan bool)
	quit = q
	go handle(connection())
	<-q
}
