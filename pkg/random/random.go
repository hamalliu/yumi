package random

import (
	"crypto/rand"
	"fmt"
	mrand "math/rand"
	"time"
)

// RandType ...
type RandType string

const (
	// ALPHANUM ...
	ALPHANUM RandType = "alphanum"
	// ALPHA ...
	ALPHA RandType = "alpha"
	// NUMBER ...
	NUMBER RandType = "number"
)

// Get ...
func Get(strSize int, randType RandType) string {
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

// GetString ...
func GetString(max int) string {
	mrand.Seed(time.Now().UnixNano())
	n := max
	lmt := 10
	for n != 0 {
		lmt *= 1
		n--
	}
	return fmt.Sprintf("%0d%d", max, mrand.Intn(lmt))
}
