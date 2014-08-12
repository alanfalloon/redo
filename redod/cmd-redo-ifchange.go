package main

import (
	"fmt"
	"path"
)

func redo_ifchange(req req, cwd string, parent target) (resp resp) {
	n := len(req.Argv) - 1
	result_chan := make(chan targetResult, n)
	for _, tgtpath := range req.Argv[1:] {
		if !path.IsAbs(tgtpath) {
			tgtpath = path.Join(cwd, tgtpath)
		}
		demand_target(tgtpath, result_chan, NEEDS_SCAN)
	}
	for ; n > 0; n-- {
		r := <-result_chan
		switch r.outcome {
		case ERROR:
			resp.ExitCode = 1
			resp.Errlines = append(resp.Errlines,
				fmt.Sprintf("failed: %s", r.name))
		case MISSING:
			resp.ExitCode = 1
			resp.Errlines = append(resp.Errlines,
				fmt.Sprintf("no rule to make target: %s", r.name))
		default:
		}
	}
	return
}
