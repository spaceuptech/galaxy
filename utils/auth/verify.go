package auth

import (
	"errors"
	"fmt"

	"github.com/dgrijalva/jwt-go"
)

// VerifyToken checks if the token is valid and returns the token claims
func (m *Module) VerifyToken(token string) (map[string]interface{}, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	// return nil, nil
	// Parse the JWT token
	tokenObj, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {

		// Declare the variables
		var alg string
		var key interface{}

		alg = jwt.SigningMethodRS256.Alg()
		key = m.config.publicKey

		// Don't forget to validate the alg is what you expect:
		if token.Method.Alg() != alg {
			return nil, errors.New("invalid signing method")
		}

		// Return the key
		return key, nil
	})
	if err != nil {
		return nil, err
	}

	// Get the claims
	if claims, ok := tokenObj.Claims.(jwt.MapClaims); ok && tokenObj.Valid {
		tokenClaims := make(map[string]interface{}, len(claims))
		for key, val := range claims {
			tokenClaims[key] = val
		}

		return tokenClaims, nil
	}

	return nil, errors.New("token could not be verified")
}

func (m *Module) VerifyCliLogin(userName, pass string) bool {
	if userName == m.config.userName && pass == m.config.key {
		return true
	}
	return false
}

func (m *Module) GenerateLoginToken() (string, error) {
	token := jwt.New(jwt.SigningMethodRS256)
	claims := make(jwt.MapClaims)
	claims["account"] = m.config.userName
	claims["role"] = "admin"
	token.Claims = claims

	tokenString, err := token.SignedString(m.config.privateKey)
	if err != nil {
		return "", fmt.Errorf("error generating token for login - %v", err)
	}
	return tokenString, nil
}
