// Package application contains the application services that orchestrate domain logic.
package application

import (
	"context"
	"fmt"
	"log"

	"github.com/jairoprogramador/fastdeploy-auth/internal/application/dto"
	"github.com/jairoprogramador/fastdeploy-auth/internal/application/mapper"
	"github.com/jairoprogramador/fastdeploy-auth/internal/domain/auth/services"
	authVos "github.com/jairoprogramador/fastdeploy-auth/internal/domain/auth/vos"
	"github.com/jairoprogramador/fastdeploy-auth/internal/domain/config/ports"
	configVos "github.com/jairoprogramador/fastdeploy-auth/internal/domain/config/vos"
)

type AuthenticateAppService struct {
	providerRepository ports.ProviderRepository
	authService        services.AuthService
	workDir            string
}

func NewAuthenticateAppService(
	providerRepository ports.ProviderRepository,
	authService services.AuthService,
	workDir string) *AuthenticateAppService {
	return &AuthenticateAppService{
		providerRepository: providerRepository,
		authService:        authService,
		workDir:            workDir,
	}
}

func (s *AuthenticateAppService) Authenticate(ctx context.Context, step, env string) (*dto.Result, error) {
	cfg, err := s.providerRepository.Load(s.workDir)
	if err != nil {
		return nil, err
	}

	log.Printf("Provider identified: %s", cfg.Name)

	var authToken *authVos.Token
	if cfg.Name != configVos.DefaultProvider {
		creds := authVos.Credentials{
			ClientID:     "my-client-id",
			ClientSecret: "my-secret",
			Scope:        "https://management.azure.com/.default",
			Extra:        map[string]string{"tenant_id": "tenant-001"},
		}

		token, err := s.authService.Authenticate(ctx, cfg.Name, creds)
		if err != nil {
			return nil, fmt.Errorf("authentication failed: %w", err)
		}
		authToken = token
	}

	return mapper.ResultToDto(authToken, step, env, s.workDir), nil
}
