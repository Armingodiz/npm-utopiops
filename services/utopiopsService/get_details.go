package utopiopsService

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
	"utopiops-cli/models"
	"utopiops-cli/utils"

	"github.com/spf13/viper"
)

func (manager *UtopiopsManager) GetApplications(token, idToken string) ([]models.AppDetail, error) {
	url := viper.GetString("CORE_URL") + "/v3/applications/environment/application"
	out, err := getList(url, token, idToken, manager.HttpHelper)
	if err != nil {
		return nil, err
	}
	var res []models.AppDetail
	if err = json.Unmarshal(out, &res); err != nil {
		return nil, err
	}
	return res, nil
}
func (manager *UtopiopsManager) GetEnvironments(token, idToken string) ([]models.EnvrionmentDetail, error) {
	url := viper.GetString("CORE_URL") + "/v3/environment"
	out, err := getList(url, token, idToken, manager.HttpHelper)
	if err != nil {
		return nil, err
	}
	var res []models.EnvrionmentDetail
	if err = json.Unmarshal(out, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func getList(url, token, idToken string, httpHelper utils.HttpHelper) ([]byte, error) {
	headers := []utils.Header{
		{
			Key:   "Authorization",
			Value: fmt.Sprintf("Bearer %s", token),
		},
		{Key: "Cookie",
			Value: fmt.Sprintf("id_token=%s", idToken),
		},
	}
	response, err, status, _ := httpHelper.HttpRequest(http.MethodGet, url, nil, headers, time.Minute, true)
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		err = errors.New("not ok with status: " + strconv.Itoa(status))
		return nil, err
	}
	return response, err
}
