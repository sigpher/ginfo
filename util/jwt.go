package util

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var mySigningKey = []byte("time@still")
var shortTokenSigningKey = []byte("short@secret")
var longTokenSigningKey = []byte("long@secret")

const TokenExpireDuration = time.Second * 24
const ShortTokenExpireDuration = time.Second * 20
const LongTokenExpireDuration = time.Second * 60

type CustomClaims struct {
	ID       uint   `json:"id" form:"id"`
	Username string `json:"username" form:"username"`
	jwt.RegisteredClaims
}

func GenToken(id uint, username string) (string, error) {
	claims := CustomClaims{
		id,
		username,

		jwt.RegisteredClaims{
			Issuer:    "GinForm",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExpireDuration)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(mySigningKey)
}

func GenDoubleToken(id uint, username string) (string, string, error) {
	sT := CustomClaims{
		id,
		username,

		jwt.RegisteredClaims{
			Issuer:    "GinForm",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ShortTokenExpireDuration)),
		},
	}

	lT := CustomClaims{
		id,
		username,

		jwt.RegisteredClaims{
			Issuer:    "GinForm",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(LongTokenExpireDuration)),
		},
	}
	shortToken := jwt.NewWithClaims(jwt.SigningMethodHS256, sT)

	shortTokenSigned, err := shortToken.SignedString(shortTokenSigningKey)
	if err != nil {
		fmt.Println("获取short token失败")
		return "", "", errors.New("")
	}
	longToken := jwt.NewWithClaims(jwt.SigningMethodHS256, lT)

	longTokenSigned, err := longToken.SignedString(longTokenSigningKey)
	if err != nil {
		fmt.Println("获取short token失败")
		return "", "", errors.New("")
	}
	return shortTokenSigned, longTokenSigned, nil
}

func ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(t *jwt.Token) (any, error) {
		return mySigningKey, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func ParseDoubleToken(sToken, lToken string) (*CustomClaims, bool, error) {
	shortToken, err := jwt.ParseWithClaims(sToken, &CustomClaims{}, func(t *jwt.Token) (any, error) {
		return shortTokenSigningKey, nil
	})
	if err != nil {
		return nil, false, errors.New("parse short token error")
	}

	if claims, ok := shortToken.Claims.(*CustomClaims); ok && shortToken.Valid {
		return claims, true, nil
	}

	longToken, err := jwt.ParseWithClaims(lToken, &CustomClaims{}, func(t *jwt.Token) (any, error) {
		return shortTokenSigningKey, nil
	})

	if claims, ok := longToken.Claims.(*CustomClaims); ok && shortToken.Valid {
		return claims, false, nil
	}

	if err != nil {
		return nil, false, errors.New("parse long token error")
	}
	return nil, false, errors.New("invalid token")
}
