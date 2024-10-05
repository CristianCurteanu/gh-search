package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTAuth struct {
	secretKey []byte
}

func NewJWTAuth(secretKey string) *JWTAuth {
	return &JWTAuth{[]byte(secretKey)}
}

func (ja *JWTAuth) CreateToken(accessToken string, expiresAt *time.Time) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": accessToken,
		"iss": "gh-search-app",
		"aud": ja.getRole(),
		"exp": expiresAt.Unix(),
		"iat": time.Now().Unix(),
	})

	tokenString, err := claims.SignedString(ja.secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (ja *JWTAuth) VerifyToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return ja.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return token, nil
}

func (ja *JWTAuth) getRole() string {
	return "user"
}
