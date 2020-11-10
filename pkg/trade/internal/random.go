package internal

import (
	"crypto/rand"
)

//RandType ...
type RandType string

const (
	//ALPHANUM ...
	ALPHANUM RandType = "alphanum"
	//ALPHA ...
	ALPHA    RandType = "alpha"
	//NUMBER ...
	NUMBER   RandType = "number"
)

//CreateRandomStr ...
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

//GetNonceStr ...
func GetNonceStr() string {
	return CreateRandomStr(30, ALPHANUM)
}
