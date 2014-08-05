package main

import (
	"github.com/alanfalloon/redo/util"
	"os"
	"path"
	"syscall"
)

func run(dofile, cwd, tgt, base string) (err error) {
	tgtpath := path.Join(cwd, tgt)
	log := logWrap("run "+tgtpath+":", log)
	tmpout := randfilename(cwd + "/outfile.")
	stdout := tmpfile(cwd + "/outstd.")
	defer os.Remove(stdout.Name())
	defer os.Remove(tmpout)
	defer stdout.Close()
	syscall.CloseOnExec(int(stdout.Fd()))
	args := []string{"-xe", dofile, tgt, base, path.Base(tmpout)}
	log.Printf("dofile=%v cwd=%v tgt=%v base=%v args=%v stdout=%v",
		dofile, cwd, tgt, base, args, stdout.Name())
	// FIXME: This is to be treated as minimal/do during tests;
	// remove once basic functionality works.
	err = os.Setenv("DO_BUILT", "true")
	check(err)
	cmd, conn := util.Launch("sh", args, cwd, stdout)
	go handle(conn, cwd)
	err = cmd.Wait()
	if err != nil {
		return
	}
	st_stdout, err := stdout.Stat()
	check(err)
	_, err = os.Lstat(tmpout)
	var has_stdout bool = st_stdout.Size() > 0
	var has_tmpout bool = err == nil
	var missing_tmpout bool = os.IsNotExist(err)
	log.Printf("Done stdout=%v tmpout=%v missing=%v", has_stdout, has_tmpout, missing_tmpout)
	switch {
	case has_stdout && has_tmpout:
		panic(tgtpath + " modified both $3 and stdout")
	case has_stdout && missing_tmpout:
		err = os.Rename(stdout.Name(), tgtpath)
	case !has_stdout && has_tmpout:
		err = os.Rename(tmpout, tgtpath)
	case !has_stdout && missing_tmpout:
		err = nil
	default:
		check(err)
	}
	return
}
