package dto

import "github.com/jairoprogramador/fastdeploy-auth/internal/domain/auth/vos"

type Result struct {
	Token   *vos.Token
	Step    string
	Env     string
	WorkDir string
}