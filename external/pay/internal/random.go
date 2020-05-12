package internal

import (
	"crypto/rand"
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

func GetNonceStr() string {
	return CreateRandomStr(30, ALPHANUM)
}
