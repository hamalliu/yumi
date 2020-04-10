package consts

import "fmt"

const (
	HeaderUser      = "yumi-user"
	HeaderTimestamp = "yumi-timestamp"
	HeaderEncrypt   = "yumi-encrypt"
)

func GetHeaders() string {
	return fmt.Sprintf("%s,%s,%s", HeaderUser, HeaderTimestamp, HeaderEncrypt)
}
