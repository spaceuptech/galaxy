package services

import (
	"fmt"
	"github.com/spaceuptech/launchpad/model"
	"github.com/spaceuptech/launchpad/runner/services/do"
	"strings"
	"sync"
)

// ManagedServices contains the map of Provider interface
type ManagedServices struct {
	lock      sync.RWMutex
	providers map[string]Provider
}

// New checks the provider
func New(config *Config) (*ManagedServices, error) {
	providers := map[string]Provider{}

	for _, provider := range config.Providers {
		array := strings.Split(provider, ":")
		isTechProvided := len(array) == 2

		switch TypeProvider(array[0]) {
		case ProviderDO:
			p, _ := do.New(config.DOToken, config.Region)
			if isTechProvided {
				providers[array[1]] = p
				continue
			}

			for _, tech := range do.GetAllTech() {
				providers[tech] = p
			}

		default:
			return nil, fmt.Errorf("invalid vendor (%s) provided", array[0])
		}
	}
	return &ManagedServices{providers: providers}, nil
}

// Provider describes the inerface a provider mush implement
type Provider interface {
	Apply(service *model.ManagedService) error
	Delete(id string) error
	GetServices() (*model.GetServiceDetails, error)
}
