package conf

import (
	"github.com/lestrrat-go/jwx/jwa"

	"yumi/pkg/types"
)

type OnlyOffice struct {
	ConfigPath        string
	SiteUrl           string
	CommandUrl        string
	ConverterUrl      string
	TempStorageUrl    string
	ApiUrl            string
	PreloaderUrl      string
	DocumentServerUrl string
	ViewedDocs        []string
	EditedDocs        []string
	ConvertedDocs     []string
	StorageFolder     string
	StoragePath       string
	SamplesPath       string
	MaxFileSize       types.SpaceSize
	MobileRegEx       string
	Static            []Static
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
