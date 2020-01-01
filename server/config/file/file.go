package file

import (
	"sync"

	"github.com/spaceuptech/launchpad/model"
)

type Manager struct {
	sync.RWMutex

	galaxyConfig *Config
	cb           func(*Config) error
}

// TODO READ CONFIG DURING INIT

func Init() *Manager {
	return &Manager{galaxyConfig: &Config{Projects: make([]*model.CreateProject, 0)}}
}
