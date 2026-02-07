package ports

import (
	"github.com/jairoprogramador/vex-client/internal/domain/logger/aggregates"
	"github.com/jairoprogramador/vex-client/internal/domain/logger/entities"
)

type Logger interface {
	Start(contextData map[string]string) *aggregates.Logger
	AddRun(logger *aggregates.Logger, stepName string) (*entities.RunRecord, error)
}
