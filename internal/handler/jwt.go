package handler

import (
	"github.com/golang-jwt/jwt"
)

func GenerateToken(payload, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"payload": payload,
	})
	return token.SignedString([]byte(secret))
}

func ParseToken(tokenString, secret string) (string, error) {
	jwtToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return "", err
	}
	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok || !jwtToken.Valid {
		return "", err
	}
	return claims["payload"].(string), nil
}
