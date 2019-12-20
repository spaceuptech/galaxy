package driver

import (
	"fmt"

	"github.com/spaceuptech/launchpad/model"
	"github.com/spaceuptech/launchpad/utils/auth"
)

// New creates a new instance of the driver module
func New(auth *auth.Module, c *Config) (Driver, error) {

	switch c.DriverType {
	case TypeIstio:
		return NewIstioDriver(auth, c)
	default:
		return nil, fmt.Errorf("invalid driver type (%s) provided", c.DriverType)
	}
}

// Config describes the configuration required by the driver module
type Config struct {
	DriverType     Type
	ConfigFilePath string
	IsInCluster    bool
	ProxyPort      uint32
}

// Driver is the interface of the modules which interact with the deployment targets
type Driver interface {
	CreateProject(project *model.Project) error
	ApplyService(service *model.Service) error
	AdjustScale(service *model.Service, activeReqs int32) error
	WaitForService(project, service string) error
	Type() Type
}

// Type is used to describe which deployment target is to be used
type Type string

const (
	// TypeIstio is the driver type used to target istio on kubernetes
	TypeIstio Type = "istio"
)
