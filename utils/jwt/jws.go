package jwt

import (
	"errors"

	"github.com/dgrijalva/jwt-go"
)

const mySigningKey = "doudoubuwangchuxinfangdeshizhong,.+-*\\"

func GetToken(user string) (string, error) {
	claims := &jwt.StandardClaims{
		ExpiresAt: 15000,
		Issuer:    user,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(mySigningKey))
}

func ValidToken(tokenstr string) (string, error) {
	token, err := jwt.Parse(tokenstr, func(token *jwt.Token) (interface{}, error) {
		return []byte(mySigningKey), nil
	})
	if err != nil {
		return "", err
	}

	if token.Valid {
		if claims, ok := token.Claims.(jwt.StandardClaims); ok {
			return claims.Issuer, nil
		} else {
			return "", errors.New("无效token")
		}
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			return "", errors.New("错误token")
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			return "", errors.New("无效token")
		} else {
			return "", err
		}
	} else {
		return "", errors.New("无效token")
	}
}
