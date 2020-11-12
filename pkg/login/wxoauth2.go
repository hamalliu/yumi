package login

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

// Scope 授权范围
type Scope string

const (
	// ScopeBase 可以获取openid
	ScopeBase Scope = "snsapi_base"
	// ScopeUserInfo 可以获取unionid
	ScopeUserInfo Scope = "snsapi_userinfo"
)

// GetWxLoginURLFromWx 从微信内（微信浏览器）登录，获取微信授权登录url
func GetWxLoginURLFromWx(appID, redirectURI, state string, scope Scope) string {
	redirectURI = url.PathEscape(redirectURI)

	vals := ""
	vals += fmt.Sprintf("appid=%s", appID)
	vals += fmt.Sprintf("&redirect_uri=%s", redirectURI)
	vals += fmt.Sprintf("&response_type=%s", "code")
	vals += fmt.Sprintf("&scope=%s", string(scope))
	vals += fmt.Sprintf("&state=%s", state)

	return fmt.Sprintf("https://open.weixin.qq.com/connect/oauth2/authorize?%s#wechat_redirect", vals)
}

// GetWxLoginURLFromBrowser 从非微信浏览器登录，获取授权登录二维码URL
func GetWxLoginURLFromBrowser(appID, redirectURI, state string) string {
	redirectURI = url.PathEscape(redirectURI)

	vals := ""
	vals += fmt.Sprintf("appid=%s", appID)
	vals += fmt.Sprintf("&redirect_uri=%s", redirectURI)
	vals += fmt.Sprintf("&response_type=%s", "code")
	vals += fmt.Sprintf("&scope=%s", "snsapi_login")
	vals += fmt.Sprintf("&state=%s", state)

	return fmt.Sprintf("https://open.weixin.qq.com/connect/qrconnect?%s#wechat_redirect", vals)
}

// TokenInfo 令牌相关信息
type TokenInfo struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`

	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"` // 秒
	RefreshToken string `json:"refresh_token"`
	OpenID       string `json:"openid"`
	Scope        string `json:"scope"`
	UnionID      string `json:"unionid"`
}

// GetToken 获取token
func GetToken(appID, appSecrt, code string) (token TokenInfo, err error) {
	vals := make(url.Values)
	vals.Add("appid", appID)
	vals.Add("secret", appSecrt)
	vals.Add("code", code)
	vals.Add("grant_type", "authorization_code")

	tokenURL := fmt.Sprintf("https://api.weixin.qq.com/sns/oauth2/access_token?%s", vals.Encode())
	resp, err := http.Get(tokenURL)
	if err != nil {
		return
	}

	err = json.NewDecoder(resp.Body).Decode(&token)
	if err != nil {
		return
	}

	if token.ErrCode != 0 {
		return token, errors.New(token.ErrMsg)
	}

	return token, nil
}

// RefreshToken 刷新token
func RefreshToken(appID, refreshToken string) (token TokenInfo, err error) {
	vals := make(url.Values)
	vals.Add("appid", appID)
	vals.Add("grant_type", "refresh_token")
	vals.Add("refresh_token", refreshToken)

	tokenURL := fmt.Sprintf("https://api.weixin.qq.com/sns/oauth2/refresh_token?%s", vals.Encode())
	resp, err := http.Get(tokenURL)
	if err != nil {
		return
	}

	err = json.NewDecoder(resp.Body).Decode(&token)
	if err != nil {
		return
	}

	if token.ErrCode != 0 {
		return token, errors.New(token.ErrMsg)
	}

	return token, nil
}

// WxUserInfo 微信用户信息
type WxUserInfo struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`

	OpenID     string   `json:"openid"`     // 用户唯一表示（对于同一公众号）
	NickName   string   `json:"nickname"`   // 昵称
	Sex        string   `json:"sex"`        // 性别
	Province   string   `json:"province"`   // 省份
	City       string   `json:"city"`       // 城市
	Country    string   `json:"country"`    // 国家
	HeadIMGURL string   `json:"headimgurl"` // 头像url
	Privilege  []string `json:"privilege"`  // 特权信息
	UnionID    string   `json:"unionid"`    // 用户统一标识。针对一个微信开放平台帐号下的应用，同一用户的unionid是唯一的。
}

// GetWxUserInfo 获取微信用户信息
func GetWxUserInfo(accessToken, openID string) (userInfo WxUserInfo, err error) {
	vals := make(url.Values)
	vals.Add("access_token", accessToken)
	vals.Add("openid", openID)

	tokenURL := fmt.Sprintf("https://api.weixin.qq.com/sns/userinfo?%s", vals.Encode())
	resp, err := http.Get(tokenURL)
	if err != nil {
		return
	}

	err = json.NewDecoder(resp.Body).Decode(&userInfo)
	if err != nil {
		return
	}

	if userInfo.ErrCode != 0 {
		return userInfo, errors.New(userInfo.ErrMsg)
	}

	return userInfo, nil
}
