package factories

import (
	"os"
	"path/filepath"

	applic "github.com/jairoprogramador/fastdeploy/internal/application"
	appPor "github.com/jairoprogramador/fastdeploy/internal/application/ports"
	"github.com/jairoprogramador/fastdeploy/internal/domain/project/ports"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/auth"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/docker"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/executor"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/path"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/project"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/variable"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/logger"
	"github.com/mattn/go-isatty"
)

type ServiceFactory interface {
	BuildOrderService() (*applic.ExecutorService, error)
	BuildInitService() (*applic.InitializeService, error)
	BuildPresenter() appPor.Presenter
}

type serviceFactory struct{}

func NewServiceFactory() ServiceFactory {
	return &serviceFactory{}
}

func (f *serviceFactory) BuildPresenter() appPor.Presenter {
	return logger.NewConsolePresenter()
}

func (f *serviceFactory) BuildInitService() (*applic.InitializeService, error) {
	projectRepository, workDir, err := f.getProjectRepository()
	if err != nil {
		return nil, err
	}

	generatorID := project.NewShaGeneratorID()

	appLogger := applic.NewAppLogger()

	inputService := project.NewSurveyUserInputService()
	return applic.NewInitializeService(filepath.Base(workDir), projectRepository, inputService, generatorID, appLogger), nil
}

func (f *serviceFactory) BuildOrderService() (*applic.ExecutorService, error) {
	projectRepository, workDir, err := f.getProjectRepository()
	if err != nil {
		return nil, err
	}

	isTerminal := isatty.IsTerminal(os.Stdout.Fd())

	pathService := path.NewPathService()

	appLogger := applic.NewAppLogger()

	cmdExecutor := executor.NewShellExecutor()
	dockerService := docker.NewDockerService(cmdExecutor)

	authService := auth.NewAuthService()
	variableResolver := variable.NewVariableResolver()

	fileConfig, err := projectRepository.Load()
	if err != nil {
		return nil, err
	}

	return applic.NewExecutorService(
		fileConfig,
		isTerminal,
		workDir,
		pathService.GetFastdeployHome(),
		projectRepository,
		dockerService,
		authService,
		variableResolver,
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
