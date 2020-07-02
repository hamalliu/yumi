package doc_service

import (
	"fmt"
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

func ConvertUriErrorMessage(errorCode int) (errorMessage string) {
	switch errorCode {
	case -20:
		errorMessage = "Error encrypt signature"
		break
	case -8:
		errorMessage = "Error document signature"
		break
	case -7:
		errorMessage = "Error document request"
		break
	case -6:
		errorMessage = "Error database"
		break
	case -5:
		errorMessage = "Error unexpected guid"
		break
	case -4:
		errorMessage = "Error download error"
		break
	case -3:
		errorMessage = "Error convertation error"
		break
	case -2:
		errorMessage = "Error convertation timeout"
		break
	case -1:
		errorMessage = "Error convertation unknown"
		break
	case 0:
		break
	default:
		errorMessage = fmt.Sprintf("%s%d", "ErrorCode = ", errorCode)
		break
	}

	return
}

func CommandServiceErrorMessage(errorCode int) (errorMessage string) {
	switch errorCode {
	case 0:
		errorMessage = "No error"
	case 1:
		errorMessage = "Document key is missing or no document with such key could be found."
	case 2:
		errorMessage = "Callback url not correct."
	case 3:
		errorMessage = "Internal server error."
	case 4:
		errorMessage = "No changes were applied to the document before the forcesave command was received."
	case 5:
		errorMessage = "Command not correct."
	case 6:
		errorMessage = "Invalid token."
	default:
		errorMessage = fmt.Sprintf("%s%d", "ErrorCode = ", errorCode)
	}

	return
}
