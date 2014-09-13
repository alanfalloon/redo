package bottomup

import (
	"os/exec"
)

type tgt interface {
	get(alias string, resp tgtobserver)
}

type tgtobserver chan<- tgtresp

type tgtresp struct {
	alias, name string
	state tgtstate
}

type tgtstate int
const (
	UPDATED tgtstate = iota
	UNCHANGED
	ERROR
	MISSING
)

type tgtreq interface {
	do(tgts map[string]tgt)
}

type tgtgetreq struct {
	path string
	resp tgtobserver
}
func (r tgtgetreq) do(tgts map[string]tgt) {
	tgt, ok := tgts[r.path]
	if !ok {
		tgt = newtgt(r.path)
		tgts[r.path] = tgt
	}
	tgt.get(r.path, r.resp)
}


var tgtinbox = func () chan<- tgtreq {
	tgtinbox := make(chan tgtreq)
	go func(tgtinbox <-chan tgtreq) {
		tgts := make(map[string]tgt)
		for req := range tgtinbox {
			req.do(tgts)
		}
	}(tgtinbox)
	return tgtinbox
}()

func gettgt(path string, resp tgtobserver) {
	tgtinbox <- tgtgetreq{path, resp}
}

type aliastgt struct {
	observers []tgtobserver
}
func (t *aliastgt) get(_ string, resp tgtobserver) {
	t.observers = append(t.observers, resp)
}

func newtgt(path string) tgt {
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

type claimtgtreq struct {
	path string
	tgt tgt
}
func (r claimtgtreq) do(tgts map[string]tgt) {
	tgts[r.path] = r.tgt
}
func claimtgt(path string, tgt tgt) {
	tgtinbox <- claimtgtreq{path, tgt}
}

type realtgt struct {
	path string
}
func newrealtgt(path string) tgt {
	return realtgt{path}
}
func (r realtgt) get(path string, resp tgtobserver) {
	resp <- tgtresp{path, r.path, ERROR}
}
