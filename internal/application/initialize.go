package application

import (
	"context"
	"runtime"

	logAgg "github.com/jairoprogramador/fastdeploy/internal/domain/logger/aggregates"
	proPor "github.com/jairoprogramador/fastdeploy/internal/domain/project/ports"
	proSer "github.com/jairoprogramador/fastdeploy/internal/domain/project/services"
	proVos "github.com/jairoprogramador/fastdeploy/internal/domain/project/vos"

	appPor "github.com/jairoprogramador/fastdeploy/internal/application/ports"
)

const MessageProjectAlreadyExists = "project already initialized, fdconfig.yaml exists"

type InitializeService struct {
	projectRepository proPor.ProjectRepository
	inputService      proPor.UserInputService
	projectName       string
	generatorID       proSer.GeneratorID
	logger            appPor.Logger
}

func NewInitializeService(
	projectName string,
	repository proPor.ProjectRepository,
	inputSvc proPor.UserInputService,
	generatorID proSer.GeneratorID,
	logger appPor.Logger) *InitializeService {
	return &InitializeService{
		projectRepository: repository,
		inputService:      inputSvc,
		projectName:       projectName,
		generatorID:       generatorID,
		logger:            logger,
	}
}

func (s *InitializeService) Run(ctx context.Context, interactive bool) (*logAgg.Logger, error) {
	logContext := map[string]string{
		"process": "initialize",
	}
	runLog := s.logger.Start(logContext)

	runRecord, err := s.logger.AddRun(runLog, "initialize")
	if err != nil {
		return runLog, err
	}

	exists, err := s.projectRepository.Exists()
	if err != nil {
		runRecord.MarkAsFailure(err)
		return runLog, err
	}
	if exists {
		runRecord.SetResult(MessageProjectAlreadyExists)
		runRecord.MarkAsWarning()
		return runLog, nil
	}

	var cfg *proVos.Config
	if interactive {
		cfg, err = s.gatherConfigFromUser(ctx)
		if err != nil {
			runRecord.MarkAsFailure(err)
			return runLog, err
		}
	} else {
		cfg = s.gatherDefaultConfig()
	}

	cfg.Project.ID = s.generatorID.ProjectID(cfg)

	err = s.projectRepository.Save(cfg)
	if err != nil {
		runRecord.MarkAsFailure(err)
		return runLog, err
	}

	runRecord.MarkAsSuccess()
	return runLog, nil
}

func (s *InitializeService) gatherConfigFromUser(ctx context.Context) (*proVos.Config, error) {
	cfg := s.gatherDefaultConfig()

	var err error

	cfg.Project.Name, err = s.inputService.Ask(ctx, "Project Name", cfg.Project.Name)
	if err != nil {
		return nil, err
	}
	cfg.Project.Version, err = s.inputService.Ask(ctx, "Project Version", cfg.Project.Version)
	if err != nil {
		return nil, err
	}
	cfg.Project.Team, err = s.inputService.Ask(ctx, "Project Team", cfg.Project.Team)
	if err != nil {
		return nil, err
	}
	cfg.Project.Organization, err = s.inputService.Ask(ctx, "Project Organization", cfg.Project.Organization)
	if err != nil {
		return nil, err
	}

	templateUrl, err := s.inputService.Ask(ctx, "Template URL", cfg.Template.URL())
	if err != nil {
		return nil, err
	}
	cfg.Template = proVos.NewTemplate(templateUrl, "")

	cfg.Runtime.Image.Source, err = s.inputService.Ask(ctx, "Runtime Image Source", cfg.Runtime.Image.Source)
	if err != nil {
		return nil, err
	}

	cfg.Runtime.Image.Tag, err = s.inputService.Ask(ctx, "Runtime Image Tag", cfg.Runtime.Image.Tag)
	if err != nil {
		return nil, err
	}

	cfg.Runtime.Volumes = s.getVolumes()

	cfg.Runtime.Env = s.getEnvVars()

	cfg.State.Backend, err = s.inputService.Ask(ctx, "State Backend", cfg.State.Backend)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func (s *InitializeService) gatherDefaultConfig() *proVos.Config {
	return &proVos.Config{
		Project: proVos.Project{
			Name:         s.projectName,
			Version:      proVos.DefaultProjectVersion,
			Team:         proVos.DefaultProjectTeam,
			Description:  proVos.DefaultProjectDescription,
			Organization: proVos.DefaultProjectOrganization,
		},
		Template: proVos.NewTemplate(proVos.DefaultUrl, proVos.DefaultRef),
		Runtime: proVos.Runtime{
			Image: proVos.Image {
				Source: proVos.DefaultImageSource,
				Tag:    proVos.DefaultImageTag,
			},
		},
		State: proVos.State{
			Backend: proVos.DefaultStateBackend,
			URL:     proVos.DefaultStateURL,
		},
	}
}

func (s *InitializeService) getVolumes() []proVos.Volume {
	var homeM2Path string
	if runtime.GOOS == "windows" {
		homeM2Path = "%USERPROFILE%\\.m2\\"
	} else {
		homeM2Path = "$HOME/.m2/"
	}

	volumes := make([]proVos.Volume, 2)
	volumes[0] = proVos.Volume{
		Host:      homeM2Path,
		Container: "/home/fastdeploy/.m2",
	}
	volumes[1] = proVos.Volume{
		Host:      "/var/run/docker.sock",
		Container: "/var/run/docker.sock",
	}
	return volumes
}

func (s *InitializeService) getEnvVars() []proVos.EnvVar {
	env := make([]proVos.EnvVar, 4)
	env[0] = proVos.EnvVar{
		Name:  "ARM_CLIENT_ID",
		Value: "{env.AZURE_CLIENT_ID}",
	}
	env[1] = proVos.EnvVar{
		Name:  "ARM_CLIENT_SECRET",
		Value: "{env.AZURE_CLIENT_SECRET}",
	}
	env[2] = proVos.EnvVar{
		Name:  "ARM_TENANT_ID",
		Value: "{env.AZURE_TENANT_ID}",
	}
	env[3] = proVos.EnvVar{
		Name:  "ARM_SUBSCRIPTION_ID",
		Value: "{env.AZURE_SUBSCRIPTION_ID}",
	}
	return env
}