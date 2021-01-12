package codec

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
)

var (
	// ErrorIVLength IV length must equal 16
	ErrorIVLength = errors.New("codec: IV length must equal 16")
)

// CbcEncrypt ...
func CbcEncrypt(key, src []byte, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	padded := pkcs7Padding(src, block.BlockSize())
	crypted := make([]byte, len(padded))

	ivl := len(iv)
	if ivl != 16 {
		return nil, ErrorIVLength
	}
	encrypter := cipher.NewCBCEncrypter(block, iv)
	encrypter.CryptBlocks(crypted, padded)

	return crypted, nil
}

// CbcDecrypt ...
func CbcDecrypt(key, src []byte, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	padded := pkcs7Padding(src, block.BlockSize())
	decrypted := make([]byte, len(padded))

	ivl := len(iv)
	if ivl != 16 {
		return nil, ErrorIVLength
	}
	decrypter := cipher.NewCBCDecrypter(block, iv)
	decrypter.CryptBlocks(decrypted, padded)

	return decrypted, nil
}

// CbcDecryptBase64 ...
func CbcDecryptBase64(key, src string, iv []byte) (string, error) {
	keyBytes, err := getKeyBytes(key)
	if err != nil {
		return "", err
	}

	encryptedBytes, err := base64.StdEncoding.DecodeString(src)
	if err != nil {
		return "", err
	}

	decryptedBytes, err := CbcDecrypt(keyBytes, encryptedBytes, iv)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(decryptedBytes), nil
}

// CbcEncryptBase64 ...
func CbcEncryptBase64(key, src string, iv []byte) (string, error) {
	keyBytes, err := getKeyBytes(key)
	if err != nil {
		return "", err
	}

	srcBytes, err := base64.StdEncoding.DecodeString(src)
	if err != nil {
		return "", err
	}

	encryptedBytes, err := CbcEncrypt(keyBytes, srcBytes, iv)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(encryptedBytes), nil
}
