package main

import (
	"fmt"
	"github.com/alanfalloon/redo/util"
	"os/exec"
	"path"
)

func fulfill_one(req req, base_cwd string) (resp resp) {
	for _, tgtpath := range req.Argv[1:] {
		cwd, tgt := path.Split(tgtpath)
		if !path.IsAbs(cwd) {
			cwd = path.Join(base_cwd, cwd)
		}
		switch err := run(tgt+".do", cwd, tgt, tgt); e := err.(type) {
		case nil:
		case *exec.ExitError:
			resp.ExitCode = 1
			resp.Errlines = append(resp.Errlines,
				fmt.Sprintf("failed: %s: %s", path.Join(cwd, tgt), e))
			return
		default:
			util.Check(err)
		}
	}
	return
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
