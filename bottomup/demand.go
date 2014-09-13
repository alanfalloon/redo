package bottomup

import (
	"github.com/alanfalloon/redo/bottomup/tgt"
)


type realtgt struct {
	path string
}
func newrealtgt(path string) tgt.T {
	return realtgt{path}
}
func (r realtgt) Watch(path string, resp tgt.Observer) {
	resp <- tgt.Update{path, r.path, tgt.ERROR}
}
