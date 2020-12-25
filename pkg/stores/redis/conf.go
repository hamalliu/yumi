package redis

import "crypto/tls"

//Config ...
type Config struct {
	Addr         string
	WriteTimeout int // second
	ReadTimeout  int // second
	DailDataBase uint8
	Password     string
	ClientName   string
	UseTLS       bool
	SkipVerify   bool
	TLSConfig    *tls.Config
}
