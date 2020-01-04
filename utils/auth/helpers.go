package auth

import "crypto/rsa"

func (m *Module) GetUserName() string {
	return m.config.userName
}

func (m *Module) GetPublicKey() *rsa.PublicKey {
	return m.config.publicKey
}
