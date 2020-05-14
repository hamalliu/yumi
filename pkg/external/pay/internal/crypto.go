package internal

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

func AesEncrypt(orig string, key []byte) string {
	// 转成字节数组
	origData := []byte(orig)
	//k := []byte(key)

	// 分组秘钥
	// NewCipher该函数限制了输入k的长度必须为16, 24或者32
	block, _ := aes.NewCipher(key)

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

	return base64.StdEncoding.EncodeToString(cryted)
}

func AesDecrypt(cryted string, key []byte) string {
	// 转成字节数组
	crytedByte, _ := base64.StdEncoding.DecodeString(cryted)
	//k := []byte(key)

	// 分组秘钥
	block, _ := aes.NewCipher(key)

	// 获取秘钥块的长度
	blockSize := block.BlockSize()

	// 加密模式
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])

	// 解密
	orig := make([]byte, len(crytedByte))
	blockMode.CryptBlocks(orig, crytedByte)

	// 去补全码
	orig = pKCS7UnPadding(orig)
	return string(orig)
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
