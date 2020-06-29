package main

import (
	"context"
	"fmt"
	"hash/crc32"
)

func Print(ctx context.Context) {
	<-ctx.Done()

	fmt.Println("hello")
}

func String(s string) int {
	v := int(crc32.ChecksumIEEE([]byte(s)))

	return v
}

func main() {
	s := `123456abcd`

	num := String(s)
	fmt.Println(num)

	//var s = []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}
	//
	//fmt.Println(Get(s))
}

func Get(tmp []string) [][]string {
	var all [][]string
	for i := range tmp {
		other := []string{}
		for _, v := range tmp {
			if v != tmp[i] {
				other = append(other, v)
			}
		}
		item := []string{}
		item = append(item, tmp[i])

		if len(other) > 0 {
			subs := Get(other)
			for _, sub := range subs {
				item = append(item, sub...)
				all = append(all, item)

				item = []string{}
				item = append(item, tmp[i])
			}
		} else {
			all = append(all, item)
		}
	}
	return all
}
