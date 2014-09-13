package bottomup

import (
	"os/exec"
	"github.com/alanfalloon/redo/bottomup/tgt"
)

type aliastgt struct {
	observers []tgt.Observer
}
func (t *aliastgt) Watch(_ string, resp tgt.Observer) {
	t.observers = append(t.observers, resp)
}

type aliasfactory struct {
	root string
	tgts tgt.TgtsI
	realtgtfac tgt.Factory
}

func (f aliasfactory) Create(path string) tgt.T {
	tgt := new(aliastgt)
	go func(tgt aliastgt) {
		// Find the canonical name to use for this target,
		// then set up all future requests to go to it, then
		// forward on all the ones we accumulated.
		cpath := canonical_name(path)
		if rel, err := filepath.Rel(f.root, cpath), err == nil && !rel.StartsWith("../") {
			cpath = rel
		}
		realtgt := f.realtgtfac.Create(cpath)
		f.tgts.Claim(cpath, realtgt)
		if path != cpath {
			f.tgts.Claim(path, realtgt)
		}
		for _, o := range tgt.observers {
			f.tgts.Watch(path, o, nil)
		}
	}(*tgt)
	return tgt
}

func canonical_name(path string) string {
	ret, err := exec.Command("readlink", "-m", path).Output()
	if err != nil {
		panic(err)
	}
	return string(ret)
}
