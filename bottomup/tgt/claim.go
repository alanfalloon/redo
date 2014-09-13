package tgt


func (t tgtthread) Claim(path string, tgt T) {
	t.inbox <- claimtgtreq{path, tgt}
}
type claimtgtreq struct {
	path string
	tgt T
}
func (r claimtgtreq) do(tgts tgts) {
	tgts[r.path] = r.tgt
}
