package bottomup

import (
	"os/exec"
	"github.com/alanfalloon/redo/bottomup/tgt"
)


func Foo() tgt.T {
	return nil
}


type aliastgt struct {
	observers []tgtobserver
}
func (t *aliastgt) get(_ string, resp tgtobserver) {
	t.observers = append(t.observers, resp)
}

func newtgt(path string) T {
	tgt := new(aliastgt)
	go func(tgt aliastgt) {
		// Find the canonical name to use for this target,
		// then set up all future requests to go to it, then
		// forward on all the ones we accumulated.
		apathb, err := exec.Command("readlink", "-m", path).Output()
		if err != nil {
			panic(err)
		}
		apath := string(apathb)
		realtgt := newrealtgt(apath)
		claimtgt(apath, realtgt)
		if path != apath {
			claimtgt(path, realtgt)
		}
		for _, o := range tgt.observers {
			gettgt(path, o)
		}
	}(*tgt)
	return tgt
}


type realtgt struct {
	path string
}
func newrealtgt(path string) T {
	return realtgt{path}
}
func (r realtgt) get(path string, resp tgtobserver) {
	resp <- tgtresp{path, r.path, tgt.ERROR}
}
