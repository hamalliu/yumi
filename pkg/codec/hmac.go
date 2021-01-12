package codec

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"io"
)

// Hmac ...
func Hmac(key []byte, body string) []byte {
	h := hmac.New(sha256.New, key)
	io.WriteString(h, body)
	return h.Sum(nil)
}

// HmacBase64 ...
func HmacBase64(key []byte, body string) string {
	return base64.StdEncoding.EncodeToString(Hmac(key, body))
}
