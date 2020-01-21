package main

import (
	"errors"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/spaceuptech/galaxy/cmd"
	"github.com/spaceuptech/galaxy/model"
	"github.com/spaceuptech/galaxy/proxy"
	"github.com/spaceuptech/galaxy/runner"
	"github.com/spaceuptech/galaxy/runner/driver"
	"github.com/spaceuptech/galaxy/runner/services"
	"github.com/spaceuptech/galaxy/server"
	"github.com/spaceuptech/galaxy/utils/auth"
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
			DriverType:     model.DriverType(driverType),
			ConfigFilePath: driverConfig,
			IsInCluster:    !outsideCluster,
		},
		Providers: &services.Config{
			DOToken:   c.String("do-token"),
			Region:    c.String("region"),
			Providers: c.StringSlice("providers"),
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

	// Set the log level
	setLogLevel(loglevel)

	s := server.New(&server.Config{Port: port})
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

func actionStartCode(c *cli.Context) error {
	envID := c.String("env")
	service, loginResp, err := cmd.CodeStart(envID)
	if err != nil {
		return err
	}
	actionCodeStruct := &model.ActionCode{
		Service:  service,
		IsDeploy: false,
	}
	if err := cmd.RunDockerFile(actionCodeStruct, loginResp); err != nil {
		return err
	}
	return nil
}

func actionBuildCode(c *cli.Context) error {
	envID := c.String("env")
	service, loginResp, err := cmd.CodeStart(envID)
	if err != nil {
		return err
	}
	actionCodeStruct := &model.ActionCode{
		Service:  service,
		IsDeploy: true,
	}
	if err := cmd.RunDockerFile(actionCodeStruct, loginResp); err != nil {
		return err
	}
	return nil
}

func actionLogin(c *cli.Context) error {
	userName := c.String("username")
	key := c.String("key")
	serverUrl := c.String("url")
	local := c.Bool("local")
	url := "url1" // todo set default url
	if local {
		url = "localhost:4122"
	}
	if serverUrl != "default url" { // todo get default url
		url = serverUrl
	}
	return cmd.LoginStart(userName, key, url, local)
}

func actionSetup(c *cli.Context) error {
	if err := cmd.CodeSetup(); err != nil {
		return err
	}
	return nil
}
