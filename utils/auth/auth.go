package auth

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"io/ioutil"
	"sync"

	"github.com/dgrijalva/jwt-go"
)

// Module manages the auth module
type Module struct {
	lock sync.RWMutex

	// For internal use
	config *Config
}

// Config is the object used to configure the auth module
type Config struct {
	// JWT related stuff
	publicKey  *rsa.PublicKey  // for RSA
	privateKey *rsa.PrivateKey // for RSA
	Secret     string          // for HSA

	// User authentication
	userName string
	key      string

	// For proxy authentication
	ProxySecret string

	Mode OperatingMode
}

// OperatingMode indicates the mode of operation
type OperatingMode string

const (
	// Runner indicates that the operating mode is runner
	Runner OperatingMode = "runner"

	// Server indicates that the operating mode is server
	Server OperatingMode = "server"
)

func Init(mode OperatingMode, username, key, jwtSecret string) *Config {
	return &Config{userName: username, key: key, Secret: jwtSecret, Mode: mode}
}

// New creates a new instance of the auth module
func New(config *Config, jwtPublicKeyPath, jwtPrivatePath string) (*Module, error) {
	m := &Module{config: config}

	// The runner needs to fetch the public key from the server for rsa
	if config.Mode == Runner {
		// Attempt fetching public key
		if success := m.fetchPublicKey(); !success {
			return nil, errors.New("could not initialise the auth module")
		}

		// Start the public key fetch routine
		go m.routineGetPublicKey()
	}
	// The server need to fetch the keys from local storage
	if config.Mode == Server {
		signBytes, err := ioutil.ReadFile(jwtPrivatePath)
		if err != nil {
			fmt.Errorf("error reading private key from path")
		}

		privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(signBytes)
		if err != nil {
			fmt.Errorf("error parsing private key")
		}

		verifyBytes, err := ioutil.ReadFile(jwtPublicKeyPath)
		if err != nil {
			fmt.Errorf("error reading public key from path")

		}

		publicKey, err := jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
		if err != nil {
			fmt.Errorf("error parsing public key")
		}

		m.config.privateKey = privateKey
		m.config.publicKey = publicKey
	}

	return m, nil
}
