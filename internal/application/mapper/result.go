package mapper

import (
	"github.com/jairoprogramador/fastdeploy-auth/internal/application/dto"
	"github.com/jairoprogramador/fastdeploy-auth/internal/domain/auth/vos"
)

func ResultToDto(token *vos.Token, step, env, workDir string) *dto.Result {
	return &dto.Result{
		Token: token,
		Step: step,
		Env: env,
		WorkDir: workDir,
	}
}