package ports

import (
	"context"

	"github.com/jairoprogramador/fastdeploy/internal/domain/docker/vos"
	domEnt "github.com/jairoprogramador/fastdeploy/internal/domain/logger/entities"
)

type DockerService interface {
	Check(ctx context.Context, stepRecord *domEnt.RunRecord) error
	Build(ctx context.Context, opts vos.BuildOptions, stepRecord *domEnt.RunRecord) error
	Run(ctx context.Context, opts vos.RunOptions, stepRecord *domEnt.RunRecord) (string, error)
}
