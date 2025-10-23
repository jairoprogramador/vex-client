package factories

import (
	"os"
	"path/filepath"
	"io"

	appPor "github.com/jairoprogramador/fastdeploy/internal/application/ports"
	applic "github.com/jairoprogramador/fastdeploy/internal/application"

	"github.com/jairoprogramador/fastdeploy/internal/domain/project/ports"

	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/docker"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/executor"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/logger"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/path"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/project"
	"github.com/mattn/go-isatty"
)

type ServiceFactory interface {
	BuildLogService(logFile io.WriteCloser) *applic.LogService
	BuildOrderService(logFile io.WriteCloser) (*applic.OrderService, error)
	BuildInitService(logFile io.WriteCloser) (*applic.InitService, error)
	BuildFileLogRepository() appPor.LogRepository
}

type serviceFactory struct {}

func NewServiceFactory() ServiceFactory {
	return &serviceFactory{}
}

func (f *serviceFactory) BuildFileLogRepository() appPor.LogRepository {
	pathService := path.NewPathService()
	return logger.NewFileLogRepository(pathService)
}

func (f *serviceFactory) BuildLogService(logFile io.WriteCloser) *applic.LogService {
	logRepo := f.BuildFileLogRepository()
	appLogger := logger.NewLoggerService(os.Stdout, logFile, false)
	logService := applic.NewLogService(logRepo, appLogger)
	return logService
}

func (f *serviceFactory) BuildInitService(logFile io.WriteCloser) (*applic.InitService, error) {
	projectRepository, workDir, err := f.getProjectRepository()
	if err != nil {
		return nil, err
	}

	appLogger := logger.NewLoggerService(os.Stdout, logFile, false)
	inputService := project.NewSurveyUserInputService()
	return applic.NewInitService(filepath.Base(workDir), projectRepository, inputService, appLogger), nil
}

func (f *serviceFactory) BuildOrderService(logFile io.WriteCloser) (*applic.OrderService, error) {
	projectRepository, workDir, err := f.getProjectRepository()
	if err != nil {
		return nil, err
	}

	isTerminal := isatty.IsTerminal(os.Stdout.Fd())

	pathService := path.NewPathService()

	appLogger := logger.NewLoggerService(os.Stdout, logFile, false)

	cmdExecutor := executor.NewShellExecutor(appLogger)
	dockerService := docker.NewDockerService(cmdExecutor, appLogger)

	return applic.NewOrderService(
		isTerminal,
		workDir,
		pathService.GetFastdeployHome(),
		projectRepository,
		dockerService,
		appLogger), nil
}

func (f *serviceFactory) getProjectRepository() (ports.ProjectRepository, string, error) {
	workDir, err := os.Getwd()
	if err != nil {
		return nil, "", err
	}
	projectRepository := project.NewYAMLProjectRepository(workDir)
	return projectRepository, workDir, nil
}

