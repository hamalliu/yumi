package security

import (
	"bytes"
	"crypto/sha256"
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

var (
	// ErrInvalidTime ...
	ErrInvalidTime = errors.New("invalid time")
	// ErrInvalidSignature ...
	ErrInvalidSignature = errors.New("invalid signature")
)

// Security ...
type Security struct {
	Key       []byte
	Timestamp string
	Nonce     string
	Signature string
}

// VerifySignature ...
func (s *Security) VerifySignature(r *http.Request, tolerance time.Duration) error {
	seconds, err := strconv.ParseInt(s.Timestamp, 10, 64)
	if err != nil {
		return err
	}

	now := time.Now().Unix()
	toleranceSeconds := int64(tolerance.Seconds())
	if seconds+toleranceSeconds < now || now+toleranceSeconds < seconds {
		return ErrInvalidTime
	}

	signContent := strings.Join([]string{
		s.Timestamp,
		s.Nonce,
		r.Method,
		r.URL.Path,
		r.URL.RawQuery,
		computeBodySignature(r),
	}, "\n")
	actualSignature := codec.HmacBase64(s.Key, signContent)

	passed := s.Signature == actualSignature
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
