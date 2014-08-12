package main

type command func(req req, cwd string, parent target) (resp resp)

var commands = map[string]command{
	"redo":          redo,
	"redo-ifchange": redo_ifchange}
