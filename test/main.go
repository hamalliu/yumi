package main

import (
	"errors"
	"fmt"
)

type Error1 struct {
	Str string
}

func (e Error1) Error() string {
	return e.Str
}

func main() {
	var err Error1
	if errors.As(err, &Error1{""}) {
		fmt.Println("hello")
	}
}
