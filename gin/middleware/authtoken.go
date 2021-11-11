package middleware

import (
	"errors"

	"github.com/dgrijalva/jwt-go"

	"yumi/gin"
	"yumi/gin/valuer"
	"yumi/pkg/token"
	"yumi/pkg/status"
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
			c.WriteJSON(nil, status.Unauthenticated())
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.WriteJSON(nil, status.Unauthenticated())
			c.Abort()
			return
		}

		for k, v := range claims {
			switch k {
			case jwtAudience, jwtExpire, jwtID, jwtIssueAt, jwtIssuer, jwtNotBefore, jwtSubject:
				// ignore the standard claims
			case string(valuer.KeyUser):
				c.Set(valuer.KeyUser, v)
			default:
				c.Set(valuer.Key(k), v)
			}
		}

		if c.Get(valuer.KeyUser) == nil {
			c.WriteJSON(nil, status.Internal().WrapError("auth token error", errors.New("the token does not contain a user name")))
			c.Abort()
			return
		}

		c.Next()
	}
}

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
