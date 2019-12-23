package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/spaceuptech/launchpad/proxy"
	"github.com/spaceuptech/launchpad/runner"
	"github.com/spaceuptech/launchpad/runner/driver"
	"github.com/spaceuptech/launchpad/server"
	"github.com/spaceuptech/launchpad/utils/auth"
)

func actionRunner(c *cli.Context) error {
	// Get runner config flags
	port := c.String("port")
	proxyPort := c.String("proxy-port")
	loglevel := c.String("log-level")

	// Get jwt config
	jwtAlgo := auth.JWTAlgorithm(c.String("jwt-algo"))
	jwtSecret := c.String("jwt-secret")
	jwtProxySecret := c.String("jwt-proxy-secret")

	// Get driver config
	driverType := c.String("driver")
	driverConfig := c.String("driver-config")
	outsideCluster := c.Bool("outside-cluster")

	// Set the log level
	setLogLevel(loglevel)

	// Create a new runner object
	r, err := runner.New(&runner.Config{
		Port:      port,
		ProxyPort: proxyPort,
		Auth: &auth.Config{
			Mode:         auth.Runner,
			JWTAlgorithm: jwtAlgo,
			Secret:       jwtSecret,
			ProxySecret:  jwtProxySecret,
		},
		Driver: &driver.Config{
			DriverType:     driver.Type(driverType),
			ConfigFilePath: driverConfig,
			IsInCluster:    !outsideCluster,
		},
	})
	if err != nil {
		logrus.Errorf("Failed to start runner - %s", err.Error())
		os.Exit(-1)
	}

	return r.Start()
}

func actionProxy(c *cli.Context) error {
	// Get all flags
	addr := c.String("addr")
	token := c.String("token")
	loglevel := c.String("log-level")

	// Set the log level
	setLogLevel(loglevel)

	// Throw an error if invalid token provided
	if len(strings.Split(token, ".")) != 3 {
		return errors.New("invalid token provided")
	}

	// Start the proxy
	p := proxy.New(addr, token)
	return p.Start()
}

func actionServer(c *cli.Context) error {
	// Get server config flags
	port := c.String("port")
	loglevel := c.String("log-level")
	authUsername := c.String("auth-username")
	authPass := c.String("auth-pass")
	jwtPublicKeyPath := c.String("jwt-public-key-path")
	jwtPrivatePath := c.String("jwt-private-key-path")

	if authUsername == "" || authPass == "" {
		fmt.Errorf("username & pass not provided")
	}

	// Set the log level
	setLogLevel(loglevel)
	a, err := auth.New(&auth.Config{
		JWTAlgorithm: auth.RSA256,
		Mode:         auth.Server,
		UserName:     authUsername,
		Pass:         authPass,
	}, jwtPublicKeyPath,jwtPrivatePath)

	if err != nil {
		fmt.Errorf("error creating an instance of auth module")
	}

	s := server.New(&server.Config{Port: port}, a)
	s.InitRoutes()
	return s.Start()
}

func setLogLevel(loglevel string) {
	switch loglevel {
	case loglevelDebug:
		logrus.SetLevel(logrus.DebugLevel)
	case loglevelInfo:
		logrus.SetLevel(logrus.InfoLevel)
	case logLevelError:
		logrus.SetLevel(logrus.ErrorLevel)
	default:
		logrus.Errorf("Invalid log level (%s) provided", loglevel)
		logrus.Infoln("Defaulting to `info` level")
		logrus.SetLevel(logrus.InfoLevel)
	}
}
