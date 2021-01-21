package conf

import (
	"github.com/lestrrat-go/jwx/jwa"

	"yumi/pkg/types"
)

//OnlyOffice 配置
type OnlyOffice struct {
	SiteURL           string
	CommandURL        string
	ConverterURL      string
	TempStorageURL    string
	APIURL            string
	PreloaderURL      string
	DocumentServerURL string
	MobileRegEx       string
	Static            []Static
	Document          Document
	Token             Token
}

//Static 配置
type Static struct {
	Name string
	Path string
}

//Token 配置
type Token struct {
	Enable                    bool
	UseForRequest             bool
	AlgorithmRequest          string
	AuthorizationHeader       string
	AuthorizationHeaderPrefix string
	Secret                    string
	ExpiresIn                 string
}

//Document 配置
type Document struct {
	StoragePath   string
	SamplesPath   string
	ConfigPath    string
	ViewedDocs    types.ArrayString
	EditedDocs    types.ArrayString
	ConvertedDocs types.ArrayString
	MaxFileSize   types.SpaceSize
}

//GetAlg 获取签名算法
func (t Token) GetAlg() jwa.SignatureAlgorithm {
	var alg jwa.SignatureAlgorithm
	switch t.AlgorithmRequest {
	case "HS256":
		alg = jwa.HS256
	case "HS384":
		alg = jwa.HS384
	case "HS512":
		alg = jwa.HS512
	default:
		alg = jwa.HS256
	}

	return alg
}
