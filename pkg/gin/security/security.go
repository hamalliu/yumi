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
	"net/url"
	"strconv"
	"strings"
	"time"

	"yumi/pkg/codec"
	"yumi/pkg/gin/header"
	"yumi/pkg/log"
)

const (
	contentTypeEncrypted = "Encrypted"
)

const (
	keyField       = "key"
	secretField    = "secret"
	signatureField = "signature"

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

// ParseContentSecurity ...
func ParseContentSecurity(decrypters map[string]codec.RsaDecrypter, r *http.Request) (
	*ContentSecurityHeader, error) {
	contentSecurity := header.ContentSecurity(r)
	attrs := parseHeader(contentSecurity)
	fingerprint := attrs[keyField]
	secret := attrs[secretField]
	signature := attrs[signatureField]

	if len(fingerprint) == 0 || len(secret) == 0 || len(signature) == 0 {
		log.Error(r.URL.Path)
		log.Error(ErrInvalidHeader)
		return nil, ErrInvalidHeader
	}

	decrypter, ok := decrypters[fingerprint]
	if !ok {
		return nil, ErrInvalidPublicKey
	}

	decryptedSecret, err := decrypter.DecryptBase64(secret)
	if err != nil {
		return nil, ErrInvalidSecret
	}

	attrs = parseHeader(string(decryptedSecret))
	base64Key := attrs[keyField]
	timestamp := attrs[timestampField]
	nonce := attrs[nonceField]
	contentType := attrs[contentTypeField]

	key, err := base64.StdEncoding.DecodeString(base64Key)
	if err != nil {
		return nil, ErrInvalidKey
	}

	return &ContentSecurityHeader{
		Key:         key,
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

	reqPath, reqQuery := getPathQuery(r)
	signContent := strings.Join([]string{
		securityHeader.Timestamp,
		securityHeader.Nonce,
		r.Method,
		reqPath,
		reqQuery,
		computeBodySignature(r),
	}, "\n")
	actualSignature := codec.HmacBase64(securityHeader.Key, signContent)

	passed := securityHeader.Signature == actualSignature
	if !passed {
		return ErrInvalidSignature
	}

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

func getPathQuery(r *http.Request) (string, string) {
	requestURI := header.RequestURI(r)
	if len(requestURI) == 0 {
		return r.URL.Path, r.URL.RawQuery
	}

	uri, err := url.Parse(requestURI)
	if err != nil {
		return r.URL.Path, r.URL.RawQuery
	}

	return uri.Path, uri.RawQuery
}
