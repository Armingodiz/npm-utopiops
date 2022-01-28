package utopiopsService

import (
	"bytes"
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

func (manager *UtopiopsManager) CreateFunction(cr models.FunctionCredentials, token, idToken string) error {
	if err := cr.IsValid(); err != nil || viper.GetString("CORE_URL") == "" {
		if err != nil {
			return err
		}
		return errors.New("invalid create body or core url is missing")
	}
	cr = cr.SetDefaults()
	intg, err := getIntegrationName(cr.Repository)
	if err != nil {
		return err
	}
	cr.IntegrationName = intg
	url := viper.GetString("CORE_URL") + "/applications/utopiops/application/function"
	json_data, err := json.Marshal(cr)
	if err != nil {
		return errors.New("bad input body")
	}
	requestBody := bytes.NewBuffer(json_data)
	Requestheaders := []utils.Header{
		{
			Key:   "Authorization",
			Value: fmt.Sprintf("Bearer %s", token),
		},
		{
			Key:   "Cookie",
			Value: fmt.Sprintf("id_token=%s", idToken),
		},
		{
			Key:   "Content-Type",
			Value: "application/json",
		},
	}
	_, err, status, _ := manager.HttpHelper.HttpRequest(http.MethodPost, url, requestBody, Requestheaders, time.Minute, true)
	if err != nil {
		return err
	}
	if status != http.StatusOK {
		return errors.New("not ok with status: " + strconv.Itoa(status))
	}
	return nil
}
