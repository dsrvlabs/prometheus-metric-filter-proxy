package config

import (
	"os"

	"github.com/dsrvlabs/prometheus-proxy/types"
	"gopkg.in/yaml.v3"
)

// Load the configuration from the file.
func Load(filename string) (*types.Config, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	config := &types.Config{}
	err = yaml.Unmarshal(content, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}


