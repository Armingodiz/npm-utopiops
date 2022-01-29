package utopiopsService

import (
	"utopiops-cli/models"
	"utopiops-cli/services/awsService"
	"utopiops-cli/utils"
)

type UtopiopsService interface {
	Deploy(cr models.DeployToUtopiopsCredentials) error
	GetApplicationDetailes(coreUrl, app, environment, token, idToken string) (models.ApplicationDetail, error)
	Register(idmUrl, idsUrl, username, password string) (string, string, error)
	CreateStaticWebsite(cr models.StaticWebsiteCredentials, token, idToken string) error
	CreateCustomDomainStaticWebsite(cr models.StaticWebsiteCredentials, token, idToken string) error
	CreateS3StaticWebsite(cr models.S3StaticWebsiteCredentials, token, idToken string) error
	CreateEcsApplication(cr models.EcsApplicationCredentials, token, idToken string) error
	CreateDockerized(cr models.DockerizedCredentials, token, idToken string) error
	CreateFunction(cr models.FunctionCredentials, token, idToken string) error
	Watch(lcr models.Log, token, idToken string) error
	Show(lcr models.Log, token, idToken string) error
	GetStaticWebsiteDomain(app, token, idToken string) (string, error)
	GetApplications(token, idToken string) ([]models.AppDetail, error)
	GetEnvironments(token, idToken string) ([]models.EnvrionmentDetail, error)
}

type UtopiopsManager struct {
	HttpHelper utils.HttpHelper
	AwsService awsService.AwsService
}

func NewService(httpHelper utils.HttpHelper, service awsService.AwsService) UtopiopsService {
	return &UtopiopsManager{
		HttpHelper: httpHelper,
		AwsService: service,
	}
}
