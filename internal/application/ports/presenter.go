package ports

import (
	"github.com/jairoprogramador/vex-client/internal/domain/logger/aggregates"
)

type Presenter interface {
	Render(log *aggregates.Logger)
}
