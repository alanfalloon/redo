package main

// #include <sys/types.h>
// #include <sys/socket.h>
import "C"

import (
	"os"
)

func socketpair(domain, typ, protocol C.int) (s0, s1 *os.File) {
	sv := [2]C.int{-1, -1}
	_, err := C.socketpair(domain, typ, protocol, &sv[0])
	check(err)
	const name = "socketpair"
	s0 = os.NewFile(uintptr(sv[0]), name)
	s1 = os.NewFile(uintptr(sv[1]), name)
	return
}
