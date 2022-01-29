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
)

func (manager *UtopiopsManager) Deploy(cr models.DeployToUtopiopsCredentials) error {
	if err := cr.IsValid(); err != nil {
		return err
	}
	containerTags := make(map[string]string)
	for _, tag := range cr.ContainerTag {
		containerTags[tag.ContainerName] = tag.ImageTag
	}
	variables := make(map[string]interface{})
	variables["container_tags"] = containerTags
	body := make(map[string]interface{})
	body["variables"] = variables
	url := fmt.Sprintf("%s/v3/applications/environment/name/%s/application/name/%s/deploy", cr.CoreUrl, cr.Environment, cr.Application)
	headers := []utils.Header{
		{
			Key:   "Authorization",
			Value: fmt.Sprintf("Bearer %s", cr.Token),
		},
		{Key: "Cookie",
			Value: fmt.Sprintf("id_token=%s", cr.IdToken),
		},
		{
			Key:   "Content-Type",
			Value: "application/json",
		},
	}
	json_data, err := json.Marshal(body)
	if err != nil {
		return errors.New("bad input body")
	}
	requestBody := bytes.NewBuffer(json_data)
	_, err, status, _ := manager.HttpHelper.HttpRequest(http.MethodPost, url, requestBody, headers, time.Minute, true)
	if err != nil {
		return err
	}
	if status != http.StatusOK {
		return errors.New("not Ok, status: " + strconv.Itoa(status))
	}
	return nil
}
