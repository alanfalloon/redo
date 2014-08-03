package main

import (
	"github.com/alanfalloon/redo/util"
	"path"
)

func fulfill_one(req req, base_cwd string) resp {
	for _, tgtpath := range req.Argv[1:] {
		log := logWrap(tgtpath+":", log)
		cwd, tgt := path.Split(tgtpath)
		cwd = base_cwd + cwd
		args := []string{"-e", tgt + ".do", tgt, tgt, "tmp"}
		log.Printf("cwd=%v tgt=%v args=%v", cwd, tgt, args)
		cmd, conn := util.Launch("sh", args, cwd)
		go handle(conn, cwd)
		err := cmd.Wait()
		util.Check(err)
	}
	return []string{}
}

func fulfill(reqs <-chan req, cwd string) <-chan resp {
	var sink = make(chan resp, 1)
	go func(sink chan<- resp) {
		defer func() { close(sink) }()
		for req := range reqs {
			sink <- fulfill_one(req, cwd)
		}
	}(sink)
	return sink
}
