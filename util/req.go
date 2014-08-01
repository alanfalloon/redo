package util

type Req struct {
	Env  map[string]string "json:env"
	Argv []string          "json:argv"
	Cwd  string            "json:cwd"
}
