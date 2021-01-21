package docservice

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/lestrrat-go/jwx/jwt"

	"yumi/conf"
)

//DocService ...
type DocService struct {
	cfg conf.Token
}

//New ...
func New(cfg conf.Token) DocService {
	return DocService{cfg: cfg}
}

//FillJwtByURL ...
func (ds DocService) FillJwtByURL(uri string, optDataObject interface{}, optIss string, optPayloadhash interface{}) (
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

//CheckJwtHeader ...
func (ds DocService) CheckJwtHeader(req *http.Request) (map[string]interface{}, error) {
	authorization := req.Header.Get(ds.cfg.AuthorizationHeader)
	if authorization != "" {
		sign := strings.TrimPrefix(authorization, ds.cfg.AuthorizationHeaderPrefix)

		token, err := jwt.ParseString(sign, jwt.WithVerify(ds.cfg.GetAlg(), ds.cfg.Secret))
		if err != nil {
			return nil, err
		}

		return token.PrivateClaims(), nil
	}

	return nil, nil
}

//GetToken ...
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

//ReadToken ...
func (ds DocService) ReadToken(tokenStr string) (map[string]interface{}, error) {
	token, err := jwt.ParseString(tokenStr, jwt.WithVerify(ds.cfg.GetAlg(), ds.cfg.Secret))
	if err != nil {
		return nil, err
	}

	return token.PrivateClaims(), nil
}

//ConvertURIErrorMessage ...
func (ds DocService) ConvertURIErrorMessage(errorCode int) (errorMessage error) {
	switch errorCode {
	case -20:
		errorMessage = fmt.Errorf("Error encrypt signature")
		break
	case -8:
		errorMessage = fmt.Errorf("Error document signature")
		break
	case -7:
		errorMessage = fmt.Errorf("Error document request")
		break
	case -6:
		errorMessage = fmt.Errorf("Error database")
		break
	case -5:
		errorMessage = fmt.Errorf("Error unexpected guid")
		break
	case -4:
		errorMessage = fmt.Errorf("Error download error")
		break
	case -3:
		errorMessage = fmt.Errorf("Error convertation error")
		break
	case -2:
		errorMessage = fmt.Errorf("Error convertation timeout")
		break
	case -1:
		errorMessage = fmt.Errorf("Error convertation unknown")
		break
	case 0:
		break
	default:
		errorMessage = fmt.Errorf("%s%d", "ErrorCode = ", errorCode)
		break
	}

	return
}

//CommandServiceErrorMessage ...
func (ds DocService) CommandServiceErrorMessage(errorCode int) (errorMessage error) {
	switch errorCode {
	case 0:
		errorMessage = fmt.Errorf("No error")
	case 1:
		errorMessage = fmt.Errorf("Document key is missing or no document with such key could be found")
	case 2:
		errorMessage = fmt.Errorf("Callback url not correct")
	case 3:
		errorMessage = fmt.Errorf("Internal server error")
	case 4:
		errorMessage = fmt.Errorf("No changes were applied to the document before the forcesave command was received")
	case 5:
		errorMessage = fmt.Errorf("Command not correct")
	case 6:
		errorMessage = fmt.Errorf("Invalid token")
	default:
		errorMessage = fmt.Errorf("%s%d", "ErrorCode = ", errorCode)
	}

	return
}

//func (ds DocService) GetPreload(preload map[string]interface{})
