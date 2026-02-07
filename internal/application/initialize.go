package application

import (
	"context"
	"errors"

	proAgg "github.com/jairoprogramador/vex-client/internal/domain/project/aggregates"
	proPor "github.com/jairoprogramador/vex-client/internal/domain/project/ports"
	proVos "github.com/jairoprogramador/vex-client/internal/domain/project/vos"
)

const MessageProjectAlreadyExists = "project already initialized, vexconfig.yaml exists"

type InitializeService struct {
	projectRepository proPor.ProjectRepository
	inputService      proPor.UserInputService
	versionService    proPor.Version
	projectName       string
}

func NewInitializeService(
	projectName string,
	repository proPor.ProjectRepository,
	inputSvc proPor.UserInputService,
	versionSvc proPor.Version,
) *InitializeService {
	return &InitializeService{
		projectRepository: repository,
		inputService:      inputSvc,
		versionService:    versionSvc,
		projectName:       projectName,
	}
}

func (s *InitializeService) Run(ctx context.Context, interactive bool) error {
	exists, err := s.projectRepository.Exists()
	if err != nil {
		return err
	}
	if exists {
		project, err := s.projectRepository.Load()
		if err != nil {
			return err
		}
		if project.IsIDDirty() {
			return s.projectRepository.Save(project)
		}
		return errors.New(MessageProjectAlreadyExists)
	}

	var project *proAgg.Project
	if interactive {
		project, err = s.createProjectFromUserInput()
		if err != nil {
			return err
		}
	} else {
		project, err = s.createDefaultProject()
		if err != nil {
			return err
		}
	}

	return s.projectRepository.Save(project)
}

func (s *InitializeService) createProjectFromUserInput() (*proAgg.Project, error) {
	name, err := s.inputService.Ask("Project Name", s.projectName)
	if err != nil {
		return nil, err
	}
	team, err := s.inputService.Ask("Project Team", proVos.DefaultProjectTeam)
	if err != nil {
		return nil, err
	}
	org, err := s.inputService.Ask("Project Organization", proVos.DefaultProjectOrganization)
	if err != nil {
		return nil, err
	}
	templateURL, err := s.inputService.Ask("Template URL", proVos.DefaultTemplateUrl)
	if err != nil {
		return nil, err
	}
	templateRef, err := s.inputService.Ask("Template Ref", proVos.DefaultTemplateRef)
	if err != nil {
		return nil, err
	}
	containerImage, err := s.inputService.Ask("Container Image", proVos.DefaultContainerImage)
	if err != nil {
		return nil, err
	}
	containerTag, err := s.inputService.Ask("Container Image Tag", proVos.DefaultContainerTag)
	if err != nil {
		return nil, err
	}

	projectData, err := proVos.NewProjectData(name, org, team, "")
	if err != nil {
		return nil, err
	}

	template, err := proVos.NewTemplate(templateURL, templateRef)
	if err != nil {
		return nil, err
	}

	container, err := proVos.NewImage(containerImage, containerTag)
	if err != nil {
		return nil, err
	}

	runtime := proVos.NewRuntime(container, []proVos.Volume{}, []proVos.EnvVar{}, []proVos.Argument{})

	projectID, err := s.getProjectID(projectData)
	if err != nil {
		return nil, err
	}

	return proAgg.NewProject(projectID, projectData, template, runtime)
}

func (s *InitializeService) createDefaultProject() (*proAgg.Project, error) {
	projectData, err := proVos.NewProjectData(
		s.projectName, proVos.DefaultProjectOrganization, proVos.DefaultProjectTeam, "")
	if err != nil {
		return nil, err
	}

	template, err := proVos.NewTemplate(proVos.DefaultTemplateUrl, proVos.DefaultTemplateRef)
	if err != nil {
		return nil, err
	}

	container, err := proVos.NewImage(proVos.DefaultContainerImage, proVos.DefaultContainerTag)
	if err != nil {
		return nil, err
	}

	runtime := proVos.NewRuntime(container, []proVos.Volume{}, []proVos.EnvVar{}, []proVos.Argument{})

	projectID, err := s.getProjectID(projectData)
	if err != nil {
		return nil, err
	}

	return proAgg.NewProject(projectID, projectData, template, runtime)
}

func (s *InitializeService) getProjectID(data proVos.ProjectData) (proVos.ProjectID, error) {
	generatedID := proVos.GenerateProjectID(data.Name(), data.Organization(), data.Team())
	return proVos.NewProjectID(generatedID.String())
}
