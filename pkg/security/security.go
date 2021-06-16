package security

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"yumi/pkg/codec"
)

const (
	contentTypeEncrypted = "Encrypted"
)

const (
	publicKeyField = "key"
	secretField    = "secret"
	signatureField = "signature"

	realSecretField  = "real_secret"
	contentTypeField = "content_type"
	timestampField   = "timestamp"
	nonceField       = "nonce"
)

var (
	// ErrInvalidContentType ...
	ErrInvalidContentType = errors.New("invalid content type")
	// ErrInvalidHeader ...
	ErrInvalidHeader = errors.New("invalid X-Content-Security header")
	// ErrInvalidKey ...
	ErrInvalidKey = errors.New("invalid key")
	// ErrInvalidPublicKey ...
	ErrInvalidPublicKey = errors.New("invalid public key")
	// ErrInvalidSecret ...
	ErrInvalidSecret = errors.New("invalid secret")
	// ErrInvalidTime ...
	ErrInvalidTime = errors.New("invalid time")
	// ErrInvalidSignature ...
	ErrInvalidSignature = errors.New("invalid signature")
)

// ContentSecurityHeader ...
type ContentSecurityHeader struct {
	Key         []byte
	Timestamp   string
	Nonce       string
	ContentType string
	Signature   string
}

// Encrypted ...
func (h *ContentSecurityHeader) Encrypted() bool {
	return h.ContentType == contentTypeEncrypted
}

func parseHeader(cs string) (attrs map[string]string) {
	// TODO:
	return nil
}

func computeBodySignature(r *http.Request) string {
	var buf bytes.Buffer
	tee := io.TeeReader(r.Body, &buf)
	sha := sha256.New()
	io.Copy(sha, tee)
	r.Body = ioutil.NopCloser(&buf)
	return fmt.Sprintf("%x", sha.Sum(nil))
}

// ParseContentSecurity ...
func ParseContentSecurity(decrypters map[string]codec.RsaDecrypter, r *http.Request) (
	*ContentSecurityHeader, error) {
	contentSecurity := r.Header.Get("X-Content-Security")
	attrs := parseHeader(contentSecurity)
	publicKey := attrs[publicKeyField]
	secret := attrs[secretField]
	signature := attrs[signatureField]

	if len(publicKey) == 0 || len(secret) == 0 || len(signature) == 0 {
		return nil, ErrInvalidHeader
	}

	decrypter, ok := decrypters[publicKey]
	if !ok {
		return nil, ErrInvalidPublicKey
	}

	decryptedSecret, err := decrypter.DecryptBase64(secret)
	if err != nil {
		return nil, ErrInvalidSecret
	}

	attrs = parseHeader(string(decryptedSecret))
	realSecret := attrs[realSecretField]
	timestamp := attrs[timestampField]
	nonce := attrs[nonceField]
	contentType := attrs[contentTypeField]

	realSecretBytes, err := base64.StdEncoding.DecodeString(realSecret)
	if err != nil {
		return nil, ErrInvalidKey
	}

	return &ContentSecurityHeader{
		Key:         realSecretBytes,
		Timestamp:   timestamp,
		Nonce:       nonce,
		ContentType: contentType,
		Signature:   signature,
	}, nil
}

// VerifySignature ...
func VerifySignature(r *http.Request, securityHeader *ContentSecurityHeader, tolerance time.Duration) error {
	seconds, err := strconv.ParseInt(securityHeader.Timestamp, 10, 64)
	if err != nil {
		return err
	}

	now := time.Now().Unix()
	toleranceSeconds := int64(tolerance.Seconds())
	if seconds+toleranceSeconds < now || now+toleranceSeconds < seconds {
		return ErrInvalidTime
	}

	signContent := strings.Join([]string{
		securityHeader.Timestamp,
		securityHeader.Nonce,
		r.Method,
		r.URL.Path,
		r.URL.RawQuery,
		computeBodySignature(r),
	}, "\n")
	actualSignature := codec.HmacBase64(securityHeader.Key, signContent)

	passed := securityHeader.Signature == actualSignature
	if !passed {
		return ErrInvalidSignature
	}

	return nil
}

func ParseBodySecurity(r *http.Request, securityHeader *ContentSecurityHeader) error {
	if r.Body == nil {
		return nil
	}
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return nil
	}
	decryptBodyBytes, err := codec.EcbDecrypt(securityHeader.Key, bodyBytes)
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(decryptBodyBytes)
	r.Body = io.NopCloser(buf)
	return nil
}
