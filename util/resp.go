package util

type Resp struct {
	Errlines []string "json:errlines"
	ExitCode int      "json:exitcode"
}
