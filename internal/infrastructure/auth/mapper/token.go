package mapper

import (
	"github.com/jairoprogramador/fastdeploy-auth/internal/domain/auth/vos"
	"github.com/jairoprogramador/fastdeploy-auth/internal/fdplugin"
)

func TokenToDomain(resp *fdplugin.AuthenticateResponse) *vos.Token {
	return &vos.Token{
		AccessToken:   resp.Token.AccessToken,
		TokenType:     resp.Token.TokenType,
		ExpiresAtUnix: resp.Token.ExpiresAtUnix,
	}
}