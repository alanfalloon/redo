package tgt

type T interface {
	Watch(alias string, resp Observer)
}
type Factory interface {
	Create(path string) T
}
type Observer chan<- Update

type Update struct {
	Alias, Name string
	State State
}

type State int
const (
	UPDATED State = iota
	UNCHANGED
	ERROR
	MISSING
)
