package main

import (
	"log"
	"net"
	"os/exec"
	"syscall"
)

func dialDaemon() (conn *net.UnixConn, cmd *exec.Cmd) {
	for {
		_conn, err := net.Dial("unix", "foo")
		if operr, ok := err.(*net.OpError); ok {
			switch operr.Err {
			case syscall.ENOENT, syscall.ECONNREFUSED:
				if cmd == nil {
					cmd = launchDaemon()
				} else {
					log.Fatal("Already launched, but still: ", err)
				}
			default:
				log.Fatal("unexpected error: ", operr)
			}
		} else if err != nil {
			log.Fatal("unexpected non-operation error:", err)
		} else {
			conn = _conn.(*net.UnixConn)
			return
		}
	}
	panic("unreachable")
}
