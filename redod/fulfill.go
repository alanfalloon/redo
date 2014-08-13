package main

import (
	"path"
)

/*
func old_fulfill(req req, base_cwd string, target target) {
	var resp resp
	for _, tgtpath := range req.Argv[1:] {
		cwd, tgt := path.Split(tgtpath)
		if !path.IsAbs(cwd) {
			cwd = path.Join(base_cwd, cwd)
		}
		dofile, cwd, tgt, base := find_dofile(cwd, tgt)
		switch dep, err := run(dofile, cwd, tgt, base); e := err.(type) {
		case nil:
			insert_dep(target, dep)
		case *exec.ExitError:
			resp.ExitCode = 1
			resp.Errlines = append(resp.Errlines,
				fmt.Sprintf("failed: %s: %s", path.Join(cwd, tgt), e))
			return
		default:
			util.Check(err)
		}
	}
	return
}
*/

func fulfill(reqs <-chan req, cwd string, parent target) <-chan resp {
	var sink = make(chan resp, 1)
	go func(sink chan<- resp) {
		defer func() { close(sink) }()
		for req := range reqs {
			cmd := path.Base(req.Argv[0])
			sink <- commands[cmd](req, cwd, parent)
		}
	}(sink)
	return sink
}
