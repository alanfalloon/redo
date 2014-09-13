package tgt

func Watch(path string, resp Observer, fac Factory) {
	inbox <- tgtgetreq{path, resp, fac}
}

type tgtgetreq struct {
	path string
	resp Observer
	fac Factory
}
func (r tgtgetreq) do(tgts tgts) {
	tgt, ok := tgts[r.path]
	if !ok {
		tgt = r.fac.Create(r.path)
		tgts[r.path] = tgt
	}
	tgt.Watch(r.path, r.resp)
}
