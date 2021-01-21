package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/dgrijalva/jwt-go"

	"yumi/gin"
	"yumi/gin/valuer"
	"yumi/pkg/token"
)

const (
	jwtAudience    = "aud"
	jwtExpire      = "exp"
	jwtID          = "jti"
	jwtIssueAt     = "iat"
	jwtIssuer      = "iss"
	jwtNotBefore   = "nbf"
	jwtSubject     = "sub"
	noDetailReason = "no detail reason"
)

// AuthToken token认证
func AuthToken(secret string, opts ...AuthorizeOption) gin.HandlerFunc {
	var authOpts AuthorizeOptions
	for _, opt := range opts {
		opt(&authOpts)
	}

	parser := token.NewTokenParser()
	return func(c *gin.Context) {
		token, err := parser.ParseToken(c.Request, secret, authOpts.PrevSecret)
		if err != nil {
			return
		}

		if !token.Valid {
			c.Writer.WriteHeader(http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.Writer.WriteHeader(http.StatusUnauthorized)
			return
		}

		for k, v := range claims {
			switch k {
			case jwtAudience, jwtExpire, jwtID, jwtIssueAt, jwtIssuer, jwtNotBefore, jwtSubject:
				// ignore the standard claims
			default:
				key := valuer.SwitchKey(k)
				if key != "" {
					// 框架固定的数据
					c.Set(key, v)
				} else {
					// 业务自定义的数据
					c.Context = context.WithValue(c.Context, k, v)
				}
			}
		}

		c.Next()
	}
}

var (
	errInvalidToken = errors.New("invalid auth token")
	errNoClaims     = errors.New("no auth params")
)

type (
	// AuthorizeOptions ...
	AuthorizeOptions struct {
		PrevSecret string
	}

	// AuthorizeOption ...
	AuthorizeOption func(opts *AuthorizeOptions)
)

// WithPrevSecret ...
func WithPrevSecret(secret string) AuthorizeOption {
	return func(opts *AuthorizeOptions) {
		opts.PrevSecret = secret
	}
}
