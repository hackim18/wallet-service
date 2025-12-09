package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

type JWTHelper struct {
	secret []byte
}

func NewJWT(v *viper.Viper) *JWTHelper {
	return &JWTHelper{
		secret: []byte(v.GetString("JWT_SECRET")),
	}
}

func (j *JWTHelper) GenerateAccessToken(id uuid.UUID, email string) (string, error) {
	if len(j.secret) == 0 {
		return "", errors.New("jwt secret is empty")
	}

	exp := time.Now().Add(60 * 24 * time.Minute).Unix()
	claims := jwt.MapClaims{
		"sub":   id.String(),
		"email": email,
		"exp":   exp,
		"iat":   time.Now().Unix(),
		"aud":   "wallet-service",
		"iss":   "wallet-service",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(j.secret)
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}

func (j *JWTHelper) DecodeAccessToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.secret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid or expired token")
}
