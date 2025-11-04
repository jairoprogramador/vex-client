package ports

import (
	"github.com/jairoprogramador/fastdeploy/internal/domain/logger/aggregates"
	//"github.com/jairoprogramador/fastdeploy/internal/domain/logger/entities"
)

type Presenter interface {
	Render(log *aggregates.Logger)
	//Header(log *aggregates.Logger)
	//ShowRun(runRecord *entities.RunRecord)
	//ShowTask(taskRecord *entities.TaskRecord, runRecord *entities.RunRecord)
	//FinalSummary(log *aggregates.Logger)
}
