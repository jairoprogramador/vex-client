// Package auth provides the gRPC plugin implementation of the auth service.
package auth

import (
	"context"
	"fmt"

	"github.com/jairoprogramador/fastdeploy-auth/internal/domain/auth/services"
	"github.com/jairoprogramador/fastdeploy-auth/internal/domain/auth/vos"
	"github.com/jairoprogramador/fastdeploy-auth/internal/fdplugin"
	"github.com/jairoprogramador/fastdeploy-auth/internal/infrastructure/auth/connector"
	"github.com/jairoprogramador/fastdeploy-auth/internal/infrastructure/auth/mapper"
)

type authService struct{}

func NewAuthService() services.AuthService {
	return &authService{}
}

func (s *authService) Authenticate(ctx context.Context, provider string, creds vos.Credentials) (*vos.Token, error) {
	conn, err := connector.NewPluginConnector(ctx, provider)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to auth plugin '%s': %w", provider, err)
	}
	defer conn.Close()

	client := fdplugin.NewAuthServiceClient(conn.GRPCConn())

	req := mapper.CredentialToRequest(creds)

	resp, err := client.Authenticate(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("gRPC authentication call failed: %w", err)
	}

	if resp.Token == nil {
		return nil, fmt.Errorf("plugin returned a nil token")
	}

	domainToken := mapper.TokenToDomain(resp)

	return domainToken, nil
}
