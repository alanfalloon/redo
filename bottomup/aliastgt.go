package bottomup

import (
	"os/exec"
	"os"
	"github.com/alanfalloon/redo/bottomup/tgt"
	"strings"
)

type aliastgt struct {
	observers []tgt.Observer
}
func (t *aliastgt) Watch(p string, resp tgt.Observer) {
	t.observers = append(t.observers, resp)
}

type aliasfactory struct {
	root string
	tgts tgt.TgtsI
	tgtfac tgt.Factory
}

func (f aliasfactory) Create(path string) tgt.T {
	tgt := new(aliastgt)
	go func(tgt *aliastgt) {
		// Find the canonical name to use for this target,
		// then set up all future requests to go to it, then
		// forward on all the ones we accumulated.
		cpath := canonical_name(path)
		cpath = strings.TrimPrefix(cpath, f.root + "/")
		realtgt := f.tgtfac.Create(cpath)
		f.tgts.Claim(cpath, realtgt)
		if path != cpath {
			f.tgts.Claim(path, realtgt)
		}
		for _, o := range tgt.observers {
			realtgt.Watch(path, o)
		}
	}(tgt)
	return tgt
}

func canonical_name(path string) string {
	ret, err := exec.Command("readlink", "-n", "-m", path).Output()
	if err != nil {
		panic(err)
	}
	return string(ret)
}

func mkaliasfac(tgts tgt.TgtsI, tgtfac tgt.Factory) aliasfactory {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return aliasfactory{
		canonical_name(string(cwd)),
		tgts,
		tgtfac,
	}
}
