package ports

import "github.com/jairoprogramador/fastdeploy-auth/internal/domain/config/vos"

type ProviderRepository interface {
	Load(workDir string) (*vos.Provider, error)
}
