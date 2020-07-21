package main

import (
	"fmt"
)

type S struct {

}

func (s S) Add() {
	fmt.Println("S")
}

type SS struct {
	S
}

func (ss SS) Add() {
	fmt.Println("SS")
}

func main() {
	var ss SS
	ss.Add()
}
