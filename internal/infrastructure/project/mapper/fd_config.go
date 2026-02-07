package mapper

import (
	"github.com/jairoprogramador/vex-client/internal/domain/project/aggregates"
	"github.com/jairoprogramador/vex-client/internal/domain/project/vos"
	"github.com/jairoprogramador/vex-client/internal/infrastructure/project/dto"
)

func ToDomainProject(configDto dto.ProjectDTO) (vos.ProjectID, vos.ProjectData, error) {
	id, err := vos.NewProjectID(configDto.ID)
	if err != nil {
		return vos.ProjectID{}, vos.ProjectData{}, err
	}
	data, err := vos.NewProjectData(
		configDto.Name,
		configDto.Organization,
		configDto.Team,
		configDto.Description)

	if err != nil {
		return vos.ProjectID{}, vos.ProjectData{}, err
	}
	return id, data, nil
}

func ToDomainRuntime(configDto dto.RuntimeDTO) (vos.Runtime, error) {
	container, err := vos.NewImage(
		configDto.Image,
		configDto.Tag)
	if err != nil {
		return vos.Runtime{}, err
	}

	volumes := make([]vos.Volume, 0, len(configDto.Run.Volumes))
	for _, dtoVol := range configDto.Run.Volumes {
		volume, err := vos.NewVolume(dtoVol.Host, dtoVol.Container)
		if err != nil {
			return vos.Runtime{}, err
		}
		volumes = append(volumes, volume)
	}

	envVars := make([]vos.EnvVar, 0, len(configDto.Run.Env))
	for _, dtoEnv := range configDto.Run.Env {
		envVar, err := vos.NewEnvVar(dtoEnv.Name, dtoEnv.Value)
		if err != nil {
			return vos.Runtime{}, err
		}
		envVars = append(envVars, envVar)
	}

	args := make([]vos.Argument, 0, len(configDto.Build.Args))
	for _, dtoArg := range configDto.Build.Args {
		arg, err := vos.NewArgument(dtoArg.Name, dtoArg.Value)
		if err != nil {
			return vos.Runtime{}, err
		}
		args = append(args, arg)
	}

	runtime := vos.NewRuntime(container, volumes, envVars, args)
	return runtime, nil
}

func ToDomain(configDto dto.FDConfigDTO) (*aggregates.Project, error) {

	id, data, err := ToDomainProject(configDto.Project)
	if err != nil {
		return nil, err
	}

	template, err := vos.NewTemplate(configDto.Template.URL, configDto.Template.Ref)
	if err != nil {
		return nil, err
	}

	runtime, err := ToDomainRuntime(configDto.Runtime)
	if err != nil {
		return nil, err
	}

	return aggregates.NewProject(id, data, template, runtime)
}

func ToRuntimeDto(runtime vos.Runtime) dto.RuntimeDTO {
	volumes := make([]dto.VolumeDTO, 0, len(runtime.Volumes()))
	for _, volume := range runtime.Volumes() {
		volumes = append(volumes, dto.VolumeDTO{
			Host:      volume.Host(),
			Container: volume.Container(),
		})
	}

	envVars := make([]dto.EnvVarDTO, 0, len(runtime.Env()))
	for _, envVar := range runtime.Env() {
		envVars = append(envVars, dto.EnvVarDTO{
			Name:  envVar.Name(),
			Value: envVar.Value(),
		})
	}

	args := make([]dto.BuildArgDTO, 0, len(runtime.Args()))
	for _, arg := range runtime.Args() {
		args = append(args, dto.BuildArgDTO{
			Name:  arg.Name(),
			Value: arg.Value(),
		})
	}

	return dto.RuntimeDTO{
		Image: runtime.Image().Image(),
		Tag:   runtime.Image().Tag(),
		Build: dto.BuildDTO{
			Args: args,
		},
		Run: dto.RunDTO{
			Volumes: volumes,
			Env:     envVars,
		},
	}
}

func ToDto(config *aggregates.Project) dto.FDConfigDTO {

	projectDto := dto.ProjectDTO{
		ID:           config.ID().String(),
		Name:         config.Data().Name(),
		Team:         config.Data().Team(),
		Description:  config.Data().Description(),
		Organization: config.Data().Organization(),
	}

	templateDto := dto.TemplateDTO{
		URL: config.Template().URL(),
		Ref: config.Template().Ref(),
	}

	runtimeDto := ToRuntimeDto(config.Runtime())

	return dto.FDConfigDTO{
		Project:  projectDto,
		Template: templateDto,
		Runtime:  runtimeDto,
	}
}
