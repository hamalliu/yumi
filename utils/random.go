package utils

import (
	"crypto/rand"
	"fmt"
	mrand "math/rand"
	"time"
)

type RandType string

const (
	ALPHANUM RandType = "alphanum"
	ALPHA    RandType = "alpha"
	NUMBER   RandType = "number"
)

func CreateRandomStr(strSize int, randType RandType) string {
	var dictionary string
	if randType == ALPHANUM {
		dictionary = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	} else if randType == ALPHA {
		dictionary = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	} else if randType == NUMBER {
		dictionary = "0123456789"
	} else {
		dictionary = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	}

	var bytes = make([]byte, strSize)
	_, _ = rand.Read(bytes)
	for k, v := range bytes {
		bytes[k] = dictionary[v%byte(len(dictionary))]
	}
	return string(bytes)
}

func RandomNumString(max int) string {
	mrand.Seed(time.Now().UnixNano())
	n := max
	lmt := 10
	for n != 0 {
		lmt *= 1
		n--
	}
	return fmt.Sprintf("%0%dd", max, mrand.Intn(lmt))
}
