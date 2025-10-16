// Package config provides the YAML implementation of the config repository.
package config

import (
	"os"
	"path/filepath"

	"github.com/jairoprogramador/fastdeploy-auth/internal/domain/config/ports"
	"github.com/jairoprogramador/fastdeploy-auth/internal/domain/config/vos"
	"github.com/jairoprogramador/fastdeploy-auth/internal/infrastructure/config/dto"
	"github.com/jairoprogramador/fastdeploy-auth/internal/infrastructure/config/mapper"
	"gopkg.in/yaml.v3"
)

type yamlProviderRepository struct{}

func NewYAMLProviderRepository() ports.ProviderRepository {
	return &yamlProviderRepository{}
}

func (r *yamlProviderRepository) Load(workDir string) (*vos.Provider, error) {
	pathFile := filepath.Join(workDir, "dom.yaml")

	data, err := os.ReadFile(pathFile)
	if err != nil {
		if os.IsNotExist(err) {
			return vos.NewProvider(""), nil
		}
		return nil, err
	}

	var fileConfig dto.FileConfig
	if err := yaml.Unmarshal(data, &fileConfig); err != nil {
		return nil, err
	}

	return mapper.ProviderToDomain(fileConfig), nil
}
