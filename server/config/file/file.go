package file

import (
	"sync"
)

type Manager struct {
	sync.RWMutex
	accountID string
}

// Init initializes file module
func Init(username string) (*Manager, error) {
	return &Manager{accountID: username}, nil
}
