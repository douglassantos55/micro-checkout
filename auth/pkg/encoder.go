package pkg

import "github.com/golang-jwt/jwt/v4"

type JWTEncoder struct{}

func NewJWTEncoder() *JWTEncoder {
	return &JWTEncoder{}
}

func (e *JWTEncoder) Encode(secret []byte, data map[string]any) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(data))
	return token.SignedString(secret)
}
