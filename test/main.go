package main

import (
	"context"
	"fmt"
	"time"
)

type Error1 struct {
	Str string
}

func (e Error1) Error() string {
	return e.Str
}

func Print(ctx context.Context) {
	<-ctx.Done()

	fmt.Println("hello")
}

func main() {
	ctx, _ := context.WithTimeout(context.Background(), time.Second)
	Print(ctx)
}
