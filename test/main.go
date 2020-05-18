package main

import "net/http"

type Error1 struct {
	Str string
}

func (e Error1) Error() string {
	return e.Str
}

func main() {
	var w http.ResponseWriter
}
