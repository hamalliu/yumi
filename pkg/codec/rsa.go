package codec

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"io/ioutil"
)

var (
	// ErrPrivateKey ...
	ErrPrivateKey = errors.New("private key error")
	// ErrPublicKey ...
	ErrPublicKey  = errors.New("failed to parse PEM block containing the public key")
	// ErrNotRsaKey ...
	ErrNotRsaKey  = errors.New("key type is not RSA")
)

type (
	// RsaDecrypter ...
	RsaDecrypter interface {
		Decrypt(input []byte) ([]byte, error)
		DecryptBase64(input string) ([]byte, error)
	}

	// RsaEncrypter ...
	RsaEncrypter interface {
		Encrypt(input []byte) ([]byte, error)
	}

	rsaBase struct {
		bytesLimit int
	}

	rsaDecrypter struct {
		rsaBase
		privateKey *rsa.PrivateKey
	}

	rsaEncrypter struct {
		rsaBase
		publicKey *rsa.PublicKey
	}
)

// NewRsaDecrypter ...
func NewRsaDecrypter(file string) (RsaDecrypter, error) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(content)
	if block == nil {
		return nil, ErrPrivateKey
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return &rsaDecrypter{
		rsaBase: rsaBase{
			bytesLimit: privateKey.N.BitLen() >> 3,
		},
		privateKey: privateKey,
	}, nil
}

// Decrypt ...
func (r *rsaDecrypter) Decrypt(input []byte) ([]byte, error) {
	return r.crypt(input, func(block []byte) ([]byte, error) {
		return rsaDecryptBlock(r.privateKey, block)
	})
}

// DecryptBase64 ...
func (r *rsaDecrypter) DecryptBase64(input string) ([]byte, error) {
	if len(input) == 0 {
		return nil, nil
	}

	base64Decoded, err := base64.StdEncoding.DecodeString(input)
	if err != nil {
		return nil, err
	}

	return r.Decrypt(base64Decoded)
}

// NewRsaEncrypter ...
func NewRsaEncrypter(key []byte) (RsaEncrypter, error) {
	block, _ := pem.Decode(key)
	if block == nil {
		return nil, ErrPublicKey
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	switch pubKey := pub.(type) {
	case *rsa.PublicKey:
		return &rsaEncrypter{
			rsaBase: rsaBase{
				// https://www.ietf.org/rfc/rfc2313.txt
				// The length of the data D shall not be more than k-11 octets, which is
				// positive since the length k of the modulus is at least 12 octets.
				bytesLimit: (pubKey.N.BitLen() >> 3) - 11,
			},
			publicKey: pubKey,
		}, nil
	default:
		return nil, ErrNotRsaKey
	}
}

// Encrypt ...
func (r *rsaEncrypter) Encrypt(input []byte) ([]byte, error) {
	return r.crypt(input, func(block []byte) ([]byte, error) {
		return rsaEncryptBlock(r.publicKey, block)
	})
}

func (r *rsaBase) crypt(input []byte, cryptFn func([]byte) ([]byte, error)) ([]byte, error) {
	var result []byte
	inputLen := len(input)

	for i := 0; i*r.bytesLimit < inputLen; i++ {
		start := r.bytesLimit * i
		var stop int
		if r.bytesLimit*(i+1) > inputLen {
			stop = inputLen
		} else {
			stop = r.bytesLimit * (i + 1)
		}
		bs, err := cryptFn(input[start:stop])
		if err != nil {
			return nil, err
		}

		result = append(result, bs...)
	}

	return result, nil
}

func rsaDecryptBlock(privateKey *rsa.PrivateKey, block []byte) ([]byte, error) {
	return rsa.DecryptPKCS1v15(rand.Reader, privateKey, block)
}

func rsaEncryptBlock(publicKey *rsa.PublicKey, msg []byte) ([]byte, error) {
	return rsa.EncryptPKCS1v15(rand.Reader, publicKey, msg)
}
