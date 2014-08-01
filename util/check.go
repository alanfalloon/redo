package util

import (
	"fmt"
)

func Check(err error) {
	if err != nil {
		panic(fmt.Sprint("check: ", err))
	}
}
