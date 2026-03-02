package main

import (
	"errors"
	"fmt"
	"os"
)

func checkFail(e error) {
	if e != nil {
		fmt.Fprintln(os.Stderr, "Boong error : ", e)
		panic(e)
	}
}

func checkFailMany(e []error) {
	if len(e) > 0 {
		checkFail(errors.Join(e...))
	}
}
