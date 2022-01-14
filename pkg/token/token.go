package token

import (
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/golang-jwt/jwt/request"
)

const claimHistoryResetDuration = time.Hour * 24

type (
	// ParseOption ...
	ParseOption func(parser *Parser)

	// Parser ...
	Parser struct {
		resetTime     time.Duration
		resetDuration time.Duration
		history       sync.Map
	}
)

// NewTokenParser ...
func NewTokenParser(opts ...ParseOption) *Parser {
	parser := &Parser{
		resetTime:     Now(),
		resetDuration: claimHistoryResetDuration,
	}

	for _, opt := range opts {
		opt(parser)
	}

	return parser
}

// ParseToken ...
func (tp *Parser) ParseToken(r *http.Request, secret, prevSecret string) (*jwt.Token, error) {
	var token *jwt.Token
	var err error

	if len(prevSecret) > 0 {
		count := tp.loadCount(secret)
		prevCount := tp.loadCount(prevSecret)

		var first, second string
		if count > prevCount {
			first = secret
			second = prevSecret
		} else {
			first = prevSecret
			second = secret
		}

		token, err = tp.doParseToken(r, first)
		if err != nil {
			token, err = tp.doParseToken(r, second)
			if err != nil {
				return nil, err
			}
			tp.incrementCount(second)
		} else {
			tp.incrementCount(first)
		}
	} else {
		token, err = tp.doParseToken(r, secret)
		if err != nil {
			return nil, err
		}
	}

	return token, nil
}

func (tp *Parser) doParseToken(r *http.Request, secret string) (*jwt.Token, error) {
	return request.ParseFromRequest(r, request.AuthorizationHeaderExtractor,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		}, request.WithParser(newParser()))
}

func (tp *Parser) incrementCount(secret string) {
	now := Now()
	if tp.resetTime+tp.resetDuration < now {
		tp.history.Range(func(key, value interface{}) bool {
			tp.history.Delete(key)
			return true
		})
	}

	value, ok := tp.history.Load(secret)
	if ok {
		atomic.AddUint64(value.(*uint64), 1)
	} else {
		var count uint64 = 1
		tp.history.Store(secret, &count)
	}
}

func (tp *Parser) loadCount(secret string) uint64 {
	value, ok := tp.history.Load(secret)
	if ok {
		return *value.(*uint64)
	}

	return 0
}

// WithResetDuration ...
func WithResetDuration(duration time.Duration) ParseOption {
	return func(parser *Parser) {
		parser.resetDuration = duration
	}
}

func newParser() *jwt.Parser {
	return &jwt.Parser{
		UseJSONNumber: true,
	}
}
