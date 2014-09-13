package tgt

type tgts map[string]T

type req interface {
	do(tgts tgts)
}

var inbox = func () chan<- req {
	inbox := make(chan req)
	go func(inbox <-chan req) {
		tgts := make(tgts)
		for req := range inbox {
			req.do(tgts)
		}
	}(inbox)
	return inbox
}()
