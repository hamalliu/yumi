package main

import (
	"fmt"
	"github.com/pkg/errors"
)

func db() error {
	return fmt.Errorf("%+v", errors.WithStack(errors.New("seqid 不存在")))
}

func srv() error {
	if err := db(); err != nil {
		return err
	}
	return nil
}

func main() {
	//err := srv()
	//fmt.Printf("%s", err.Error())
	fmt.Println(fmt.Sprintf(`%%%s%%`, "sdf"))
}
