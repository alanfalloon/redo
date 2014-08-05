package main

import (
	"crypto/rand"
	"encoding/base64"
	"os"
)

func tmpfile(prefix string) (f *os.File) {
collision:
	for {
		fn := randfilename(prefix)
		var err error
		f, err = os.OpenFile(fn, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0644)
		if os.IsExist(err) {
			continue collision
		}
		check(err)
		return
	}
}

func randfilename(prefix string) string {
	var b [6]byte
	_, err := rand.Read(b[:])
	check(err)
	return prefix + base64.URLEncoding.EncodeToString(b[:])
}
