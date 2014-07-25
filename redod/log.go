package main

import (
	"fmt"
	_log "log"
	"os"
)

var log *_log.Logger = _log.New(os.Stderr, fmt.Sprint("redod(", os.Getpid(), "): "), _log.LstdFlags)

func logWrap(prefix string, l *_log.Logger) *_log.Logger {
	return _log.New(os.Stderr, fmt.Sprint(l.Prefix(), prefix), l.Flags())
}
