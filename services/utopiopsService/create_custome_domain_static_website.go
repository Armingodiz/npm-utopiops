package utopiopsService

import (
	"errors"
	"utopiops-cli/models"

	"github.com/spf13/viper"
)

func (manager *UtopiopsManager) CreateCustomDomainStaticWebsite(cr models.StaticWebsiteCredentials, token, idToken string) error {
	if cr.Domain == "" {
		return errors.New("we need domain for this action")
	}
	cr = cr.SetDefaults()
	url := viper.GetString("DM_URL") + "/utopiops/static-website/setup"
	return createWithLog(cr, url, token, idToken, manager.HttpHelper)
}
