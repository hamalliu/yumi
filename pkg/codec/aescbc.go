package codec

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
)

var (
	// ErrorIVLength IV length must equal 16
	ErrorIVLength = errors.New("codec: IV length must equal 16")
)

// CbcEncrypt ...
func CbcEncrypt(key, src []byte, IV []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	padded := pkcs7Padding(src, block.BlockSize())
	crypted := make([]byte, len(padded))

	ivl := len(IV)
	if ivl != 16 {
		return nil, ErrorIVLength
	}
	encrypter := cipher.NewCBCEncrypter(block, IV)
	encrypter.CryptBlocks(crypted, padded)

	return crypted, nil
}

// CbcDecrypt ...
func CbcDecrypt(key, src []byte, IV []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	padded := pkcs7Padding(src, block.BlockSize())
	decrypted := make([]byte, len(padded))

	ivl := len(IV)
	if ivl != 16 {
		return nil, ErrorIVLength
	}
	decrypter := cipher.NewCBCDecrypter(block, IV)
	decrypter.CryptBlocks(decrypted, padded)

	return decrypted, nil
}
