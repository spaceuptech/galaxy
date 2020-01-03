package auth

import "crypto/rsa"

func (m *Module) GetUserName() string {
	return m.config.UserName
}

func (m *Module) GetPublicKey() *rsa.PublicKey {
	return m.config.PublicKey
}
