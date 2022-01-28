package models

import "errors"

type DeployToUtopiopsCredentials struct {
	CoreUrl      string
	Environment  string
	Application  string
	Token        string
	IdToken      string
	ContainerTag []ContainerTag
}

type ApplicationDetail struct {
	EcrRegisteryUrl string
	ContainerNames  []string
}

type ContainerTag struct {
	ContainerName string
	ImageTag      string
}

func (cr DeployToUtopiopsCredentials) IsValid() error {
	if cr.Application == "" {
		return errors.New("application name needed")
	}
	if cr.CoreUrl == "" {
		return errors.New("core url needed")
	}
	if cr.Environment == "" {
		return errors.New("environment name needed")
	}
	if cr.Token == "" || cr.IdToken == "" {
		return errors.New("authorization credentials are needed")
	}
	if cr.ContainerTag[0].ImageTag == "" || cr.ContainerTag[0].ContainerName == "" {
		return errors.New("container detailes are needed")
	}
	return nil
}
