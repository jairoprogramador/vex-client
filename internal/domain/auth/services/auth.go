package services

import (
	"context"
	"github.com/jairoprogramador/fastdeploy-auth/internal/domain/auth/vos"
)

type AuthService interface {
	Authenticate(ctx context.Context, provider string, creds vos.Credentials) (*vos.Token, error)
}