package application

import (
	appPor "github.com/jairoprogramador/vex-client/internal/application/ports"
	"github.com/jairoprogramador/vex-client/internal/domain/logger/aggregates"
	"github.com/jairoprogramador/vex-client/internal/domain/logger/entities"
)

type AppLogger struct{}

func NewAppLogger() appPor.Logger {
	return &AppLogger{}
}

func (l *AppLogger) Start(contextData map[string]string) *aggregates.Logger {
	return aggregates.NewLogger(contextData)
}

func (l *AppLogger) AddRun(logger *aggregates.Logger, stepName string) (*entities.RunRecord, error) {
	runRecord, err := entities.NewRunRecord(stepName)
	if err != nil {
		return nil, err
	}
	logger.AddRun(runRecord)

	return runRecord, nil
}
