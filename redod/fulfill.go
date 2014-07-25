package main

func fulfill(reqs <-chan Req) <-chan resp {
	var sink = make(chan resp, 1)
	go func(sink chan<- resp) {
		defer func() { close(sink) }()

		log := logWrap("fulfill:", log)
		log.Print("begin", reqs)
		defer log.Print("done", reqs)

		for req := range reqs {
			log.Print("handling", req)
			sink <- []string{"ok"}
		}
	}(sink)
	return sink
}
