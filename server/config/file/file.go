package file

import (
	"sync"
)

type Manager struct {
	sync.RWMutex
	accountID string
}

// TODO READ CONFIG DURING INIT

func Init(username string) (*Manager, error) {
	return &Manager{accountID: username}, nil
}
