package utopiopsService

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
	"utopiops-cli/utils"

	"github.com/spf13/viper"
)

func (manager *UtopiopsManager) GetStaticWebsiteDomain(app, token, idToken string) (string, error) {
	url := fmt.Sprintf("%s/applications/utopiops/name/%s", viper.GetString("CORE_URL"), app)
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
		return "", err
	}
	if status != http.StatusOK {
		return "", errors.New("not ok with status: " + strconv.Itoa(status))
	}
	var res struct {
		Domain string `json:"domain"`
	}
	if err = json.Unmarshal(response, &res); err != nil {
		return "", err
	}
	return res.Domain, nil
}
