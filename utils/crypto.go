package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"strings"
)

// Key ...
const Key = "buwangchuxinfangdeshizhong,.+-*\\"

// GetKey ...
func GetKey(key string, size int) string {
	for {
		if len(key) < size {
			key += key
		} else {
			break
		}
	}
	key = key[0:size]
	return key
}

// AesEncrypt ...
func AesEncrypt(orig string, key []byte) (string, error) {
	// 转成字节数组
	origData := []byte(orig)
	//k := []byte(key)

	// 分组秘钥
	// NewCipher该函数限制了输入k的长度必须为16, 24或者32
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// 获取秘钥块的长度
	blockSize := block.BlockSize()

	// 补全码
	origData = pKCS7Padding(origData, blockSize)

	// 加密模式
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])

	// 创建数组
	cryted := make([]byte, len(origData))
	
	// 加密
	blockMode.CryptBlocks(cryted, origData)

	return base64.StdEncoding.EncodeToString(cryted), nil
}

// AesDecrypt ...
func AesDecrypt(cryted string, key []byte) (string, error) {
	// 转成字节数组
	crytedByte, err := base64.StdEncoding.DecodeString(cryted)
	if err != nil {
		return "", err
	}

	// 分组秘钥
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// 获取秘钥块的长度
	blockSize := block.BlockSize()

	// 加密模式
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])

	// 解密
	orig := make([]byte, len(crytedByte))
	blockMode.CryptBlocks(orig, crytedByte)

	// 去补全码
	orig = pKCS7UnPadding(orig)
	return string(orig), nil
}

//补码
//AES加密数据块分组长度必须为128bit(byte[16])，
//密钥长度可以是128bit(byte[16])、192bit(byte[24])、256bit(byte[32])中的任意一个。
func pKCS7Padding(ciphertext []byte, blocksize int) []byte {
	padding := blocksize - len(ciphertext)%blocksize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

//去码
func pKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

// DesEncrypt ...
func DesEncrypt(origData, key []byte) (string, error) {
	block, err := des.NewCipher(key[:8])
	if err != nil {
		return "", err
	}
	origData = PKCS5Padding(origData, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, key[:8])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	encodeString := base64.StdEncoding.EncodeToString(crypted)
	return encodeString, nil
}

// DesDecrypt ...
func DesDecrypt(encodeString string, key []byte) (string, error) {
	crypted, err := base64.StdEncoding.DecodeString(encodeString)
	block, err := des.NewCipher(key[:8])
	if err != nil {
		return "", err
	}
	blockMode := cipher.NewCBCDecrypter(block, key[:8])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS5UnPadding(origData)
	origDataStr := string(origData)
	return origDataStr, nil
}

// TripleDesEncrypt ...
func TripleDesEncrypt(origData, key []byte) (string, error) {
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return "", err
	}
	origData = PKCS5Padding(origData, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, key[:8])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	encodeString := base64.StdEncoding.EncodeToString(crypted)
	return encodeString, nil
}

// TripleDesDecrypt ...
func TripleDesDecrypt(encodeString string, key []byte) (string, error) {
	defer func() {
		recover()
	}()
	crypted, err := base64.StdEncoding.DecodeString(encodeString)
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return "", err
	}
	blockMode := cipher.NewCBCDecrypter(block, key[:8])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS5UnPadding(origData)
	origDataStr := string(origData)
	return origDataStr, nil
}

// PKCS5Padding ...
func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// PKCS5UnPadding ...
func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

// MD5 ...
func MD5(data []byte) []byte {
	md5Ctx := md5.New()
	md5Ctx.Write(data)
	return md5Ctx.Sum(nil)
}

// MD5String ...
func MD5String(data []byte) string {
	return hex.EncodeToString(MD5(data))
}

// Md5LowerString ...
func Md5LowerString(data []byte) string {
	ret := strings.ToLower(hex.EncodeToString(MD5(data)))
	return ret
}

// Sha1 ...
func Sha1(data []byte) []byte {
	h := sha1.New()
	h.Write(data)
	return h.Sum(nil)
}

// Sha1String ...
func Sha1String(data []byte) string {
	ret := hex.EncodeToString(Sha1(data))
	return ret
}

// Sha1LowerString ...
func Sha1LowerString(data []byte) string {
	ret := strings.ToLower(hex.EncodeToString(Sha1(data)))
	return ret
}
