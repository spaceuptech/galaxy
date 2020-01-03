package file

import (
	"sync"
)

type Manager struct {
	sync.RWMutex
	accountID string
}

// TODO READ CONFIG DURING INIT

func Init() (*Manager, error) {
	return &Manager{}, nil
}
