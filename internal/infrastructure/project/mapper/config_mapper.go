package mapper

import (
	"github.com/jairoprogramador/fastdeploy/internal/domain/project/vos"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/project/dto"
)

func ToDomain(configDto dto.FileConfig) *vos.Config {

	if configDto.State.Backend == "" {
		configDto.State.Backend = vos.DefaultStateBackend
	}

	domainVolumes := make([]vos.Volume, 0, len(configDto.Runtime.Volumes))
	for _, dtoVol := range configDto.Runtime.Volumes {
		domainVolumes = append(domainVolumes, vos.Volume{
			Host:      dtoVol.Host,
			Container: dtoVol.Container,
		})
	}

	extraParams := make(map[string]string)
	for _, item := range configDto.Auth.Params.Extra {
		for key, value := range item {
			extraParams[key] = value
		}
	}

	domainEnvVars := make([]vos.EnvVar, 0, len(configDto.Runtime.Env))
	for _, dtoEnv := range configDto.Runtime.Env {
		domainEnvVars = append(domainEnvVars, vos.EnvVar{
			Name:  dtoEnv.Name,
			Value: dtoEnv.Value,
		})
	}

	return &vos.Config{
		Project: vos.Project{
			ID:           configDto.Project.ID,
			Name:         configDto.Project.Name,
			Version:      configDto.Project.Version,
			Team:         configDto.Project.Team,
			Description:  configDto.Project.Description,
			Organization: configDto.Project.Organization,
		},
		Template: vos.NewTemplate(configDto.Template.URL, configDto.Template.Ref),
		Runtime: vos.Runtime{
			Image: vos.Image{
				Source:      configDto.Runtime.Image.Source,
				Tag:         configDto.Runtime.Image.Tag,
				CoreVersion: configDto.Runtime.Image.CoreVersion,
			},
			Volumes: domainVolumes,
			Env:     domainEnvVars,
		},
		State: vos.State{
			Backend: configDto.State.Backend,
			URL:     configDto.State.URL,
		},
		Auth: vos.Auth{
			Plugin: configDto.Auth.Plugin,
			Params: vos.AuthParams{
				ClientID:     configDto.Auth.Params.ClientID,
				GrantType:    configDto.Auth.Params.GrantType,
				ClientSecret: configDto.Auth.Params.ClientSecret,
				Scope:        configDto.Auth.Params.Scope,
				Extra:        extraParams,
			},
		},
	}
}

func ToDto(config *vos.Config) dto.FileConfig {

	dtoVolumes := make([]dto.VolumeDTO, 0, len(config.Runtime.Volumes))
	for _, vol := range config.Runtime.Volumes {
		dtoVolumes = append(dtoVolumes, dto.VolumeDTO{
			Host:      vol.Host,
			Container: vol.Container,
		})
	}

	dtoExtraParams := make([]map[string]string, 0, len(config.Auth.Params.Extra))
	for key, value := range config.Auth.Params.Extra {
		dtoExtraParams = append(dtoExtraParams, map[string]string{key: value})
	}

	dtoEnvVars := make([]dto.EnvVarDTO, 0, len(config.Runtime.Env))
	for _, envVar := range config.Runtime.Env {
		dtoEnvVars = append(dtoEnvVars, dto.EnvVarDTO{
			Name:  envVar.Name,
			Value: envVar.Value,
		})
	}

	dtoConfig := dto.FileConfig{
		Project: dto.ProjectDTO{
			ID:           config.Project.ID,
			Name:         config.Project.Name,
			Version:      config.Project.Version,
			Team:         config.Project.Team,
			Description:  config.Project.Description,
			Organization: config.Project.Organization,
		},
		Template: dto.TemplateDTO{
			URL: config.Template.URL(),
			Ref: config.Template.Ref(),
		},
		Runtime: dto.RuntimeDTO{
			Image: dto.ImageDTO{
				Source: config.Runtime.Image.Source,
				Tag:    config.Runtime.Image.Tag,
				CoreVersion: config.Runtime.Image.CoreVersion,
			},
			Volumes: dtoVolumes,
			Env:     dtoEnvVars,
		},
		State: dto.StateDTO{
			Backend: config.State.Backend,
			URL:     config.State.URL,
		},
		Auth: dto.AuthDTO{
			Plugin: config.Auth.Plugin,
			Params: dto.AuthParamsDTO{
				ClientID:     config.Auth.Params.ClientID,
				GrantType:    config.Auth.Params.GrantType,
				ClientSecret: config.Auth.Params.ClientSecret,
				Extra:        dtoExtraParams,
			},
		},
	}
	return dtoConfig
}
