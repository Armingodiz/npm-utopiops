package awsService

import (
	"utopiops-cli/models"
)

type AwsService interface {
	Show(lcr models.Log, cr models.ProviderCredentials) error
	Watch(lcr models.Log, cr models.ProviderCredentials) error
}
type AwsManager struct {
}

func NewService() AwsService {
	return &AwsManager{}
}
