package file

import (
	"context"
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"

	"github.com/spaceuptech/launchpad/model"
	"github.com/spaceuptech/launchpad/utils"
)

func (m *Manager) StoreConfigToFile(ctx context.Context, config *Config) error {
	m.Lock()
	defer m.Unlock()

	m.galaxyConfig = config

	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("error yaml marshalling while storing config - %v", err)
	}

	return ioutil.WriteFile(utils.ConfigFilePath, data, 0644)
}

// remove removes specified element in slice
// NOTE: IT DOESN'T PRESERVE ORDER
func remove(arr interface{}, index int) (interface{}, error) {
	switch s := arr.(type) {
	case []*model.Project:
		s[index] = s[len(s)-1]
		return s[:len(s)-1], nil
	case []*model.Environment:
		s[index] = s[len(s)-1]
		return s[:len(s)-1], nil
	case []*model.Cluster:
		s[index] = s[len(s)-1]
		return s[:len(s)-1], nil
	default:
		return nil, fmt.Errorf("error removing from specified index incorrect type")
	}
}
