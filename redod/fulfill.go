package main

func fulfill_one(req Req) resp {
	return []string{"ok"}
}

func fulfill(reqs <-chan Req) <-chan resp {
	var sink = make(chan resp, 1)
	go func(sink chan<- resp) {
		defer func() { close(sink) }()
		for req := range reqs {
			sink <- fulfill_one(req)
		}
	}(sink)
	return sink
}
