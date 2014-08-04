package main

import (
	"github.com/alanfalloon/redo/util"
	"os"
	"path"
)

func run(dofile, cwd, tgt, base string) (err error) {
	tgtpath := path.Join(cwd, tgt)
	log := logWrap("run "+tgtpath+":", log)
	args := []string{"-e", dofile, tgt, base, "tmp"}
	log.Printf("dofile=%v cwd=%v tgt=%v base=%v args=%v", dofile, cwd, tgt, base, args)
	out, err := os.Create("tmp2")
	util.Check(err)
	cmd, conn := util.Launch("sh", args, cwd, out)
	go handle(conn, cwd)
	err = cmd.Wait()
	return
}
