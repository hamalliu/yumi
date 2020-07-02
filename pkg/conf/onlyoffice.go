package conf

import (
	"github.com/lestrrat-go/jwx/jwa"

	"yumi/pkg/types"
)

type OnlyOffice struct {
	SiteUrl           string
	CommandUrl        string
	ConverterUrl      string
	TempStorageUrl    string
	ApiUrl            string
	PreloaderUrl      string
	DocumentServerUrl string
	MobileRegEx       string
	Static            []Static
	Document          Document
	Token             Token
}

type Static struct {
	Name string
	Path string
}

type Token struct {
	Enable                    bool
	UseForRequest             bool
	AlgorithmRequest          string
	AuthorizationHeader       string
	AuthorizationHeaderPrefix string
	Secret                    string
	ExpiresIn                 string
}

type Document struct {
	StoragePath   string
	SamplesPath   string
	ConfigPath    string
	ViewedDocs    types.ArrayString
	EditedDocs    types.ArrayString
	ConvertedDocs types.ArrayString
	MaxFileSize   types.SpaceSize
}

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
