package mapper

import (
	"github.com/jairoprogramador/fastdeploy-auth/internal/domain/auth/vos"
	"github.com/jairoprogramador/fastdeploy-auth/internal/fdplugin"
)

func CredentialToRequest(credential vos.Credentials) *fdplugin.AuthenticateRequest {
	return &fdplugin.AuthenticateRequest{
		Config: &fdplugin.AuthConfig{
			ClientId:     credential.ClientID,
			ClientSecret: credential.ClientSecret,
			Scope:        credential.Scope,
			Extra:        credential.Extra,
		},
	}
}