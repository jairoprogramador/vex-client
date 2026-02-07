package ports

import (
	"context"

	"github.com/jairoprogramador/vex-client/internal/fdplugin"
)

type AuthService interface {
	Authenticate(ctx context.Context, provider string, request *fdplugin.AuthenticateRequest) (*fdplugin.AuthenticateResponse, error)
}
