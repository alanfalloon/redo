package main

import (
	"fmt"
)

func fulfill_one(req Req) resp {
	log := logWrap(fmt.Sprintf("fulfill[%v]:", req.Argv), log)
	log.Print("begin")
	defer log.Print("done")
	return []string{"ok"}
}

func fulfill(reqs <-chan Req) <-chan resp {
	var sink = make(chan resp, 1)
	go func(sink chan<- resp) {
		defer func() { close(sink) }()

		log := logWrap("fulfill:", log)
		log.Print("begin", reqs)
		defer log.Print("done", reqs)

		for req := range reqs {
			sink <- fulfill_one(req)
		}
	}(sink)
	return sink
}
