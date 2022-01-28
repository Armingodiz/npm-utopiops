package models

import "errors"

type PushCredentials struct {
	Repository string
	EcrUrl     string
	ImageTag   string
}

func (cr PushCredentials) IsValid() error {
	if cr.EcrUrl == "" {
		return errors.New("ecr url needed")
	}
	if cr.Repository == "" {
		return errors.New("repository neede")
	}
	if cr.ImageTag == "" {
		return errors.New("image tag id needed")
	}
	return nil
}
