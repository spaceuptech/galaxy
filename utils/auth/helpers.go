package auth

import (
	"github.com/dgrijalva/jwt-go"
)

func (m *Module) GetUserName() string {
	return m.config.userName
}

func (m *Module) GetPublicKey() string {
	return m.config.base64PublicKey
}

func (m *Module) GenerateHS256Token() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"account": m.config.userName,
		"role":    "admin",
	})
	return token.SignedString([]byte(m.config.Secret))
}
