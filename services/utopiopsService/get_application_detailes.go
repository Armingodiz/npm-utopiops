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
)

func (manager *UtopiopsManager) GetApplicationDetailes(coreUrl, app, environment, token, idToken string) (out models.ApplicationDetail, err error) {
	url := fmt.Sprintf("%s/v3/applications/environment/name/%s/application/name/%s/tf", coreUrl, environment, app)
	headers := []utils.Header{
		{
			Key:   "Authorization",
			Value: fmt.Sprintf("Bearer %s", token),
		},
		{Key: "Cookie",
			Value: fmt.Sprintf("id_token=%s", idToken),
		},
	}
	response, err, status, _ := manager.HttpHelper.HttpRequest(http.MethodGet, url, nil, headers, time.Minute, true)
	if err != nil {
		return
	}
	if status != http.StatusOK {
		err = errors.New("not ok with status: " + strconv.Itoa(status))
		return
	}
	var res struct {
		Ecr        string                   `json:"ecrRegisteryUrl"`
		Containers []map[string]interface{} `json:"containers"`
	}
	if err = json.Unmarshal(response, &res); err != nil {
		return
	}
	names := make([]string, 0)
	for _, container := range res.Containers {
		names = append(names, container["name"].(string))
	}
	return models.ApplicationDetail{
		EcrRegisteryUrl: res.Ecr,
		ContainerNames:  names,
	}, nil
}
