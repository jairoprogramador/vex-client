package application

import (
	"context"

	proPor "github.com/jairoprogramador/fastdeploy/internal/domain/project/ports"
	proSer "github.com/jairoprogramador/fastdeploy/internal/domain/project/services"
	proVos "github.com/jairoprogramador/fastdeploy/internal/domain/project/vos"
	logAgg "github.com/jairoprogramador/fastdeploy/internal/domain/logger/aggregates"

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

	cfg.Runtime.CoreVersion, err = s.inputService.Ask(ctx, "Runtime Core Version", cfg.Runtime.CoreVersion)
	if err != nil {
		return nil, err
	}
	cfg.Runtime.Image.Source, err = s.inputService.Ask(ctx, "Runtime Image Source", cfg.Runtime.Image.Source)
	if err != nil {
		return nil, err
	}

	cfg.Runtime.Image.Tag, err = s.inputService.Ask(ctx, "Runtime Image Tag", cfg.Runtime.Image.Tag)
	if err != nil {
		return nil, err
	}
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
			CoreVersion: proVos.DefaultCoreVersion,
			Image: proVos.Image{
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
