package main

func redo(req req, cwd string, parent target) (resp resp) {
	return redo_ifchange(req, cwd, parent)
}
