package runner

import (
	"github.com/spaceuptech/galaxy/runner/driver"
	"github.com/spaceuptech/galaxy/utils/auth"
)

// Config is the object required to configure the runner
type Config struct {
	Port      string
	ProxyPort string

	// Configuration for the driver
	Driver *driver.Config

	// Configuration for the auth module
	Auth *auth.Config
}
