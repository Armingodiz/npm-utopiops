package utopiopsService

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
	"utopiops-cli/models"
	"utopiops-cli/utils"

	"github.com/spf13/viper"
)

func (manager *UtopiopsManager) Watch(lcr models.Log, token, idToken string) error {
	logGroup, region, err := manager.getLogGroup(lcr.App, lcr.Environment, token, idToken)
	if err != nil {
		return err
	}
	keyId, secretKey, err := manager.getProviderDetailes(lcr.Environment, token, idToken)
	if err != nil {
		return err
	}
	cr := models.ProviderCredentials{
		AccessKeyId:     keyId,
		SecretAccessKey: secretKey,
		Region:          region,
	}
	lcr.LogGroup = logGroup
	manager.AwsService.Watch(lcr, cr)
	return nil
}

func (manager *UtopiopsManager) Show(lcr models.Log, token, idToken string) error {
	logGroup, region, err := manager.getLogGroup(lcr.App, lcr.Environment, token, idToken)
	if err != nil {
		return err
	}
	keyId, secretKey, err := manager.getProviderDetailes(lcr.Environment, token, idToken)
	if err != nil {
		return err
	}
	cr := models.ProviderCredentials{
		AccessKeyId:     keyId,
		SecretAccessKey: secretKey,
		Region:          region,
	}
	lcr.LogGroup = logGroup
	manager.AwsService.Show(lcr, cr)
	return nil
}

func (manager *UtopiopsManager) getLogGroup(app, env, token, idToken string) (logGroup, region string, err error) {
	url := fmt.Sprintf("%s/v3/applications/environment/name/%s/application/name/%s/resources", viper.GetString("CORE_URL"), env, app)
	Requestheaders := []utils.Header{
		{
			Key:   "Authorization",
			Value: fmt.Sprintf("Bearer %s", token),
		},
		{Key: "Cookie",
			Value: fmt.Sprintf("id_token=%s", idToken),
		},
	}
	response, err, status, _ := manager.HttpHelper.HttpRequest(http.MethodGet, url, nil, Requestheaders, time.Minute, true)
	if err != nil {
		return
	}
	if status != http.StatusOK {
		err = errors.New("not ok with status: " + strconv.Itoa(status))
		return
	}
	var res struct {
		ClusterName map[string]interface{} `json:"cluster_name"`
		LogGroup    dto                    `json:"log_groups"`
		Service     dto                    `json:"service"`
	}
	if err = json.Unmarshal(response, &res); err != nil {
		return
	}
	cluster := res.ClusterName["value"].(string)
	if gr, ok := res.LogGroup.Value[cluster].(string); ok {
		logGroup = gr
	} else {
		err = errors.New("log group not found")
	}
	if cl, ok := res.Service.Value["cluster"].(string); ok {
		region = strings.Split(cl, ":")[3]
	} else {
		err = errors.New("provider region not found")
	}
	return
}

type dto struct {
	Value map[string]interface{} `json:"value"`
	Type  []interface{}          `json:"type"`
}

func (manager *UtopiopsManager) getProviderDetailes(env, token, idToken string) (keyId, secretKey string, err error) {
	var respDto struct {
		Credentials struct {
			AccessKeyId     string `json:"accessKeyId"`
			SecretAccessKey string `json:"secretAccessKey"`
		}
	}
	// Get the access token and url
	url := fmt.Sprintf("%s/v3/environment/name/%s/provider/credentials", viper.GetString("CORE_URL"), env)
	Requestheaders := []utils.Header{
		{
			Key:   "Authorization",
			Value: fmt.Sprintf("Bearer %s", token),
		},
		{Key: "Cookie",
			Value: fmt.Sprintf("id_token=%s", idToken),
		},
	}
	out, err, statusCode, _ := manager.HttpHelper.HttpRequest(http.MethodGet, url, nil, Requestheaders, time.Minute, true)
	if err != nil {
		return
	}
	if statusCode != http.StatusOK {
		err = errors.New("not ok with status: " + strconv.Itoa(statusCode))
		return
	}
	err = json.Unmarshal(out, &respDto)
	if err != nil {
		return
	}
	keyId = respDto.Credentials.AccessKeyId
	secretKey = respDto.Credentials.SecretAccessKey
	return
}
