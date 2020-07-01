package doc_service

import (
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/lestrrat-go/jwx/jwt"

	"yumi/pkg/conf"
)

type DocService struct {
	cfg conf.Token
}

func New(cfg conf.Token) DocService {
	return DocService{cfg: cfg}
}

func (ds DocService) FillJwtByUrl(uri string, optDataObject interface{}, optIss string, optPayloadhash interface{}) (
	[]byte, error) {

	expireIn, err := time.ParseDuration(ds.cfg.ExpiresIn)
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

	return jwt.Sign(t, ds.cfg.GetAlg(), ds.cfg.Secret)
}

func (ds DocService) CheckJwtHeader(req *http.Response) (map[string]interface{}, error) {

	authorization := req.Header.Get(ds.cfg.AuthorizationHeader)
	if authorization != "" {
		sign := strings.TrimPrefix(authorization, ds.cfg.AuthorizationHeaderPrefix)

		token, err := jwt.ParseString(sign, jwt.WithVerify(ds.cfg.GetAlg(), ds.cfg.Secret))
		if err != nil {
			return nil, err
		}

		return token.PrivateClaims(), nil
	} else {
		return nil, nil
	}
}

func (ds DocService) GetToken(privateClaims map[string]interface{}) ([]byte, error) {

	expireIn, err := time.ParseDuration(ds.cfg.ExpiresIn)
	if err != nil {
		panic(err)
	}

	t := jwt.New()
	_ = t.Set(jwt.ExpirationKey, time.Now().Add(expireIn).Unix())
	for k, v := range privateClaims {
		_ = t.Set(k, v)
	}

	return jwt.Sign(t, ds.cfg.GetAlg(), ds.cfg.Secret)
}

func (ds DocService) ReadToken(tokenStr string) (map[string]interface{}, error) {
	token, err := jwt.ParseString(tokenStr, jwt.WithVerify(ds.cfg.GetAlg(), ds.cfg.Secret))
	if err != nil {
		return nil, err
	}

	return token.PrivateClaims(), nil
}
