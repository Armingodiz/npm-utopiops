package utopiopsService

import (
	"utopiops-cli/models"

	"github.com/spf13/viper"
)

func (manager *UtopiopsManager) CreateEcsApplication(cr models.EcsApplicationCredentials, token, idToken string) error {
	cr = cr.SetDefaults()
	url := viper.GetString("DM_URL") + "/flash-setup/stage-2"
	return createWithLog(cr, url, token, idToken, manager.HttpHelper)
}
