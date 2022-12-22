package pkg

import (
	"errors"
	"time"
)

type Auth interface {
	Authenticate(Credentials) (string, error)
}

type Encoder interface {
	Encode(secret []byte, data map[string]any) (string, error)
}

type auth struct {
	encoder Encoder
}

func NewAuth(encoder Encoder) Auth {
	return &auth{encoder: encoder}
}

func (a *auth) Authenticate(credentials Credentials) (string, error) {
	if credentials.Username != "admin" || credentials.Password != "admin" {
		return "", errors.New("invalid credentials")
	}

	claims := map[string]any{
		"sub": "YWRtaW4K",
		"iat": time.Now().Unix(),
	}

	return a.encoder.Encode([]byte("your-256-bit-secret"), claims)
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
