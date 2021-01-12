package codec

import (
	"bytes"	
	"encoding/base64"
)


func pkcs7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func pkcs7Unpadding(src []byte, blockSize int) ([]byte, error) {
	length := len(src)
	unpadding := int(src[length-1])
	if unpadding >= length || unpadding > blockSize {
		return nil, ErrPaddingSize
	}

	return src[:length-unpadding], nil
}

func getKeyBytes(key string) ([]byte, error) {
	if len(key) > 32 {
		keyBytes, err := base64.StdEncoding.DecodeString(key)
		if err != nil {
			return nil, err
		}
		return keyBytes, nil
	}

	return []byte(key), nil
}
