package mapper

import (
	"github.com/jairoprogramador/fastdeploy-auth/internal/domain/config/vos"
	"github.com/jairoprogramador/fastdeploy-auth/internal/infrastructure/config/dto"
)

func ProviderToDomain(fileConfig dto.FileConfig) *vos.Provider {
	if fileConfig.Technology.Provider == "" {
		return vos.NewProvider(vos.DefaultProvider)
	}
	return vos.NewProvider(fileConfig.Technology.Provider)
}
