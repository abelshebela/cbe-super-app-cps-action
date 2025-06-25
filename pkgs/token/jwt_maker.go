package token

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type JwtMaker struct {
	secretKey string
}

func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < 32 {
		return nil, fmt.Errorf("secret key must be at least 32 characters long")
	}

	maker := &JwtMaker{
		secretKey: secretKey,
	}

	return maker, nil
}

func (maker *JwtMaker) VerifyUserToken(tokenString string) (*Payload, error) {
	payload := &Payload{}
	token, err := jwt.ParseWithClaims(tokenString, payload, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(maker.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Payload); ok && token.Valid {
		return claims, nil
	} else {
		return nil, fmt.Errorf("invalid token")
	}
}
