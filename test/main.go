package main

import (
	"context"
	"fmt"
)

func Print(ctx context.Context) {
	<-ctx.Done()

	fmt.Println("hello")
}

func main() {

}
