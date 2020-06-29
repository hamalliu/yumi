package helplers

import (
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/lestrrat-go/jwx/jwt"

	"yumi/pkg/conf"
)

type DocumentService struct {
	UserIp string
}

func (ds DocumentService) FillJwtByUrl(uri string, optDataObject interface{}, optIss string, optPayloadhash interface{}) (
	[]byte, error) {
	cfgToken := conf.Get().Office.Token

	expireIn, err := time.ParseDuration(cfgToken.ExpiresIn)
	if err != nil {
		panic(err)
	}

	t := jwt.New()
	_ = t.Set(jwt.ExpirationKey, time.Now().Add(expireIn).Unix())
	_ = t.Set(jwt.IssuerKey, optIss)
	urlObj, _ := url.Parse(uri)
	_ = t.Set("query", urlObj.Query())
	_ = t.Set("payload", optDataObject)
	_ = t.Set("payloadhash", optPayloadhash)

	return jwt.Sign(t, cfgToken.GetAlg(), cfgToken.Secret)
}

func (ds DocumentService) CheckJwtHeader(req *http.Response) (map[string]interface{}, error) {
	cfgToken := conf.Get().Office.Token

	authorization := req.Header.Get(cfgToken.AuthorizationHeader)
	if authorization != "" {
		sign := strings.TrimPrefix(authorization, cfgToken.AuthorizationHeaderPrefix)

		cfgToken := conf.Get().Office.Token

		token, err := jwt.ParseString(sign, jwt.WithVerify(cfgToken.GetAlg(), cfgToken.Secret))
		if err != nil {
			return nil, err
		}

		return token.PrivateClaims(), nil
	} else {
		return nil, nil
	}
}

func (ds DocumentService) GetToken(privateClaims map[string]interface{}) ([]byte, error) {
	cfgToken := conf.Get().Office.Token

	expireIn, err := time.ParseDuration(cfgToken.ExpiresIn)
	if err != nil {
		panic(err)
	}

	t := jwt.New()
	_ = t.Set(jwt.ExpirationKey, time.Now().Add(expireIn).Unix())
	for k, v := range privateClaims {
		_ = t.Set(k, v)
	}

	return jwt.Sign(t, cfgToken.GetAlg(), cfgToken.Secret)
}

func (ds DocumentService) ReadToken(tokenStr string) (map[string]interface{}, error) {
	cfgToken := conf.Get().Office.Token

	token, err := jwt.ParseString(tokenStr, jwt.WithVerify(cfgToken.GetAlg(), cfgToken.Secret))
	if err != nil {
		return nil, err
	}

	return token.PrivateClaims(), nil
}
