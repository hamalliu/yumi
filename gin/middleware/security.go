package middleware

import (
	"encoding/base64"
	"errors"
	"strings"
	"time"

	"yumi/gin"
	"yumi/gin/header"
	"yumi/gin/security"
	"yumi/gin/valuer"
	"yumi/pkg/codec"
	"yumi/pkg/status"
)

var (
	// ErrInvalidHeader ...
	ErrInvalidHeader = errors.New("invalid x-yumi-content-security header")
	// ErrInvalidKey ...
	ErrInvalidKey = errors.New("invalid key")
	// ErrInvalidPublicKey ...
	ErrInvalidPublicKey = errors.New("invalid public key")
	// ErrInvalidSecret ...
	ErrInvalidSecret = errors.New("invalid secret")
)

var (
	publicKeyField = "key"
	secretField    = "secret"
	signatureField = "signature"

	realSecretField = "real_secret"
	timestampField  = "timestamp"
	nonceField      = "nonce"
)

func parseHeader(cs string) (attrs map[string]string) {
	attrs = make(map[string]string)
	as := strings.Split(cs, ";")
	for _, v := range as {
		pair := strings.Split(v, "=")
		if len(pair) != 2 {
			continue
		}
		attrs[pair[0]] = pair[1]
	}

	return attrs
}

// LoginSecurity 登录接口安全传输
func LoginSecurity(decrypters map[string]codec.RsaDecrypter, tolerance time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		contentSecurity := header.ContentSecurity(c.Request)
		attrs := parseHeader(contentSecurity)
		publicKey := attrs[publicKeyField]
		secret := attrs[secretField]
		signature := attrs[signatureField]

		if len(publicKey) == 0 || len(secret) == 0 || len(signature) == 0 {
			c.WriteJSON(nil, status.InvalidArgument().WithDetails(ErrInvalidHeader))
			c.Abort()
			return
		}

		decrypter, ok := decrypters[publicKey]
		if !ok {
			c.WriteJSON(nil, status.InvalidArgument().WithDetails(ErrInvalidPublicKey))
			c.Abort()
			return
		}

		decryptedSecret, err := decrypter.DecryptBase64(secret)
		if err != nil {
			c.WriteJSON(nil, status.InvalidArgument().WithDetails(ErrInvalidSecret))
			c.Abort()
			return
		}

		attrs = parseHeader(string(decryptedSecret))
		realSecret := attrs[realSecretField]
		timestamp := attrs[timestampField]
		nonce := attrs[nonceField]
		if len(realSecret) == 0 || len(timestamp) == 0 || len(nonce) == 0 {
			c.WriteJSON(nil, status.InvalidArgument().WithDetails(ErrInvalidHeader))
			c.Abort()
			return
		}

		realSecretBytes, err := base64.StdEncoding.DecodeString(realSecret)
		if err != nil {
			c.WriteJSON(nil, status.InvalidArgument().WithDetails(ErrInvalidKey))
			c.Abort()
			return
		}

		s := security.Security{
			Key:       realSecretBytes,
			Timestamp: timestamp,
			Nonce:     nonce,
			Signature: signature,
		}

		err = s.VerifySignature(c.Request, tolerance)
		if err != nil {
			c.WriteJSON(nil, status.InvalidArgument().WithDetails(err))
			c.Abort()
			return
		}

		c.Set(valuer.KeySecret, s.Key)
		c.Set(valuer.KeyTimestamp, s.Timestamp)
		c.Set(valuer.KeyNonce, s.Nonce)
		c.Set(valuer.KeySignature, s.Signature)
		c.Next()
	}
}

// NoLoginSecurity 非登录接口安全传输
func NoLoginSecurity(getSecret func (user string) string, tolerance time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		contentSecurity := header.ContentSecurity(c.Request)
		attrs := parseHeader(contentSecurity)
		signature := attrs[signatureField]
		timestamp := attrs[timestampField]
		nonce := attrs[nonceField]
		if len(signature) == 0 || len(timestamp) == 0 || len(nonce) == 0 {
			c.WriteJSON(nil, status.InvalidArgument().WithDetails(ErrInvalidHeader))
			c.Abort()
			return
		}

		user := c.Get(valuer.KeyUser).String()
		if user == "" {
			c.WriteJSON(nil, status.Internal().WithDetails(errors.New("security middleware: no found user name")))
			c.Abort()
			return
		}
		secret := []byte(getSecret(user))
		s := security.Security{
			Key:       secret,
			Timestamp: timestamp,
			Nonce:     nonce,
			Signature: signature,
		}

		err := s.VerifySignature(c.Request, tolerance)
		if err != nil {
			c.WriteJSON(nil, status.Unauthenticated().WithDetails(err))
			c.Abort()
			return
		}

		c.Set(valuer.KeySecret, secret)
		c.Set(valuer.KeyTimestamp, s.Timestamp)
		c.Set(valuer.KeyNonce, s.Nonce)
		c.Set(valuer.KeySignature, s.Signature)
		c.Next()
	}
}
