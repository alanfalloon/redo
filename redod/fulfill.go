package main

import (
	"fmt"
	"github.com/alanfalloon/redo/util"
	"os"
	"os/exec"
	"path"
)

func fulfill_one(req req, base_cwd string) (resp resp) {
	for _, tgtpath := range req.Argv[1:] {
		log := logWrap(tgtpath+":", log)
		cwd, tgt := path.Split(tgtpath)
		if !path.IsAbs(cwd) {
			cwd = path.Join(base_cwd, cwd)
		}
		args := []string{"-e", tgt + ".do", tgt, tgt, "tmp"}
		log.Printf("cwd=%v tgt=%v args=%v", cwd, tgt, args)
		out, err := os.Create("tmp2")
		util.Check(err)
		cmd, conn := util.Launch("sh", args, cwd, out)
		go handle(conn, cwd)
		switch err = cmd.Wait(); e := err.(type) {
		case nil:
		case *exec.ExitError:
			resp.ExitCode = 1
			resp.Errlines = append(resp.Errlines,
				fmt.Sprintf("failed: %s: %s", path.Join(cwd, tgt), e))
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
