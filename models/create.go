package models

import "errors"

type CreateCredentials interface {
	IsValid() error
}

type StaticWebsiteCredentials struct {
	Name           string `json:"name"`
	RepositoryUrl  string `json:"repositoryUrl"`
	Description    string `json:"description"`
	Index_document string `json:"index_document"`
	Error_document string `json:"error_document"`
	BuildCommand   string `json:"buildCommand"`
	OutputPath     string `json:"outputPath"`
	Branch         string `json:"branch"`
	Domain         string `json:"domainName"`
}

func (cr StaticWebsiteCredentials) IsValid() error {
	if cr.Name == "" {
		return errors.New("application name needed")
	}
	if cr.RepositoryUrl == "" {
		return errors.New("repository url needed")
	}
	if cr.BuildCommand == "" {
		return errors.New("build command needed")
	}
	if cr.OutputPath == "" {
		return errors.New("output path needed")
	}
	return nil
}

func (cr StaticWebsiteCredentials) SetDefaults() StaticWebsiteCredentials {
	if cr.Index_document == "" {
		cr.Index_document = "index.html"
	}
	if cr.Error_document == "" {
		cr.Error_document = "error.html"
	}
	if cr.Branch == "" {
		cr.Branch = "main"
	}
	if cr.Description == "" {
		cr.Description = "created from cli"
	}
	return cr
}

type S3StaticWebsiteCredentials struct {
	Applications []S3Application `json:"applications"`
	Environments []Environment   `json:"environments"`
	//	Region       string          `json:"region"`
}
type S3Application struct {
	Name                 string `json:"name"`
	Description          string `json:"description"`
	RepositoryUrl        string `json:"repositoryUrl"`
	IsDynamicApplication bool   `json:"isDynamicApplication"`
	Type                 string `json:"type"`
}
type Environment struct {
	Name string `json:"name"`
	//Arn  string `json:"arn"`
}

func (cr S3StaticWebsiteCredentials) IsValid() error {
	if cr.Applications[0].Name == "" {
		return errors.New("applicaion name needed")
	}
	if cr.Applications[0].Type == "" {
		return errors.New("applicaion type needed")
	}
	if cr.Applications[0].RepositoryUrl == "" {
		return errors.New("applicaion repo url needed")
	}
	if cr.Environments[0].Name == "" {
		return errors.New("environment name needed")
	}
	return nil
}

func (cr S3StaticWebsiteCredentials) SetDefaults() S3StaticWebsiteCredentials {
	return cr
}

type EcsApplicationCredentials struct {
	Applications []EcsApplication `json:"applications"`
	Environments []Environment    `json:"environments"`
	//Region       string           `json:"region"`
}

type EcsApplication struct {
	Name                 string `json:"name"`
	Description          string `json:"description"`
	Port                 string `json:"port"`
	Protocol             string `json:"protocol"`
	IsDynamicApplication bool   `json:"isDynamicApplication"`
	Type                 string `json:"type"`
}

func (cr EcsApplicationCredentials) IsValid() error {
	if cr.Applications[0].Name == "" {
		return errors.New("applicaion name needed")
	}
	if cr.Applications[0].Type == "" {
		return errors.New("applicaion type needed")
	}
	if cr.Applications[0].Port == "" {
		return errors.New("applicaion port needed")
	}
	if cr.Applications[0].Protocol == "" {
		return errors.New("application protocol needed")
	}
	if cr.Environments[0].Name == "" {
		return errors.New("environment name needed")
	}
	return nil
}

func (cr EcsApplicationCredentials) SetDefaults() EcsApplicationCredentials {
	cr.Applications[0].Protocol = "https"
	return cr
}

type DockerizedCredentials struct {
	Name                 string        `json:"name"`
	Description          string        `json:"description"`
	Port                 string        `json:"port"`
	Cpu                  float64       `json:"cpu"`
	Memory               int           `json:"memory"`
	EnvironmentVariables []EnvVariable `json:"environmentVariables"`
	HealthCheckPath      string        `json:"healthCheckPath"`
	Branch               string        `json:"branch"`
	Repository           string        `json:"repositoryUrl"`
	IntegrationName      string        `json:"integration_name"`
	DomainName           string        `json:"domainName"`
}

type EnvVariable struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func (cr DockerizedCredentials) IsValid() error {
	if cr.Name == "" {
		return errors.New("applicaion name needed")
	}
	if cr.Port == "" {
		return errors.New("applicaion port needed")
	}
	if cr.Branch == "" {
		return errors.New("branch needed")
	}
	if cr.Repository == "" {
		return errors.New("repo url needed")
	}
	if cr.DomainName != "" {
		if cr.Cpu > 2 {
			return errors.New("cpu must be less than 2")
		}
		if cr.Memory > 2048 {
			return errors.New("memory must be less than 2048")
		}
	}
	return nil
}

func (cr DockerizedCredentials) SetDefaults() DockerizedCredentials {
	cr.EnvironmentVariables = make([]EnvVariable, 0)
	cr.HealthCheckPath = "/health"
	if cr.DomainName == "" {
		cr.Cpu = 0.25
		cr.Memory = 128
	}
	if cr.Description == "" {
		cr.Description = "created from cli"
	}
	return cr
}

type FunctionCredentials struct {
	Name            string `json:"name"`
	Description     string `json:"description"`
	Branch          string `json:"branch"`
	Repository      string `json:"repositoryUrl"`
	IntegrationName string `json:"integration_name"`
	DomainName      string `json:"domainName"`
}

func (cr FunctionCredentials) IsValid() error {
	if cr.Name == "" {
		return errors.New("applicaion name needed")
	}
	if cr.Branch == "" {
		return errors.New("branch needed")
	}
	if cr.Repository == "" {
		return errors.New("repo url needed")
	}
	return nil
}

func (cr FunctionCredentials) SetDefaults() FunctionCredentials {
	if cr.Description == "" {
		cr.Description = "created from cli"
	}
	return cr
}
