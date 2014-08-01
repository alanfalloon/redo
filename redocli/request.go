package main

import (
	"github.com/alanfalloon/redo/util"
	"log"
	"os"
	"strings"
)

func request() util.Req {
	return util.Req{env(), os.Args, cwd()}
}

func env() map[string]string {
	env := make(map[string]string)
	for _, e := range os.Environ() {
		kv := strings.SplitN(e, "=", 2)
		k, v := kv[0], kv[1]
		if len(k) > 5 && k[0:5] == "REDO_" {
			env[k] = v
		}
	}
	return env
}

func cwd() string {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal("getcwd: ", err)
	}
	return cwd
}
