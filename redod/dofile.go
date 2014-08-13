package main

import (
	"github.com/alanfalloon/redo/util"
	"os"
	"path"
	"strings"
)

func find_dofile(tgtpath string) (dofile, cwd, tgt, base string) {
	cwd, tgt = path.Split(tgtpath)
	dofile = tgt + ".do"
	base = tgt
	if exist(cwd, dofile) {
		return
	}
	tgtparts := strings.Split(tgt, ".")
	subdir := ""
	for {
		for part := 1; part < len(tgtparts); part++ {
			dofile = "default." + strings.Join(tgtparts[part:], ".") + ".do"
			base = path.Join(subdir, strings.Join(tgtparts[:part], "."))
			if exist(cwd, dofile) {
				return
			}
		}
		if cwd == "" {
			panic(tgt)
		}
		var last string
		cwd, last = path.Split(cwd)
		subdir = path.Join(last, subdir)
		tgt = path.Join(last, tgt)
	}
}

func exist(dir, filename string) bool {
	filename = path.Join(dir, filename)
	_, err := os.Lstat(filename)
	if err == nil {
		return true
	}
	if !os.IsNotExist(err) {
		util.Check(err)
	}
	return false
}
