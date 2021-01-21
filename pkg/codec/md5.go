package codec

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
)

// Md5 ...
func Md5(cnt []byte) string {
	md5ctx := md5.New()
	md5ctx.Write(cnt)
	return strings.ToUpper(hex.EncodeToString(md5ctx.Sum(nil)))
}
