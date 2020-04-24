package main

import (
	"errors"
	"fmt"
	"time"

	"yumi/utils"
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

	fmt.Println(utils.CreateRandomStr(32, utils.ALPHANUM))
	fmt.Println(time.Now().Format("06121545.999"))
}
