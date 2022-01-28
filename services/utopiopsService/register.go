package utopiopsService

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
	"utopiops-cli/utils"
)

func (manager *UtopiopsManager) Register(idsUrl, idmUrl, username, password string) (token string, idToken string, err error) {
	loginChallenge, cookie, err := getLoginChalleng(idsUrl, manager.HttpHelper)
	if err != nil {
		return
	}
	url, err := getUrl(idmUrl, username, password, loginChallenge, manager.HttpHelper)
	if err != nil {
		return
	}
	location1, cookies, err := getLocation(url, cookie, manager.HttpHelper)
	if err != nil {
		return
	}
	location2, _, err := getLocation(location1, cookie, manager.HttpHelper)
	if err != nil {
		return
	}
	Requestheaders := []utils.Header{
		{
			Key:   "Cookie",
			Value: fmt.Sprintf("%s;%s", cookies[0], cookies[1]),
		},
	}
	_, err, status, headers := manager.HttpHelper.HttpRequest(http.MethodGet, location2, nil, Requestheaders, time.Minute, false)
	if err != nil || status != http.StatusFound {
		var errString string
		if err != nil {
			errString = err.Error()
		}
		return "", "", errors.New("error in last request, err= " + errString + " status:" + strconv.Itoa(status))
	}
	locationHeader := strings.Split(headers.Get("location"), "#")[1]
	params := strings.Split(locationHeader, "&")
	token = getParam("access_token=", params)
	idToken = getParam("id_token=", params)
	if token == "" || idToken == "" {
		err = errors.New("token or idToken was not found")
	}
	return
}
func getLoginChalleng(idsUrl string, httpHelper utils.HttpHelper) (string, string, error) {
	var url string

	if strings.Contains(idsUrl, "staging") {
		url = idsUrl + "/oauth2/auth?audience=&client_id=portal&prompt=&redirect_uri=https://portal.staging.utopiops.com/auth/login/accept&response_type=token+id_token&scope=offline+openid&state=" + utils.GetString(16) + "&nonce=" + utils.GetString(16)
	} else {
		url = idsUrl + "/oauth2/auth?audience=&client_id=portal-7n8gyf13ezsq9p4t6eixyoimmpx5s1u5&nonce=9TefwKFqMC4KnQrb&prompt=&redirect_uri=https://water.utopiops.com/auth/login/accept&response_type=token+id_token&scope=offline+openid&state=" + utils.GetString(16)
	}

	_, err, status, headers := httpHelper.HttpRequest(http.MethodGet, url, nil, nil, time.Minute, false)
	if err != nil || status != http.StatusFound {
		var errString string
		if err != nil {
			errString = err.Error()
		}
		return "", "", errors.New("error in first request, err= " + errString + " status:" + strconv.Itoa(status))
	}
	loginChallenge := strings.Replace(strings.Split(headers.Get("location"), "login_challenge")[1], "=", "", 1)
	return loginChallenge, headers.Get("Set-Cookie"), nil
}
func getUrl(idmUrl, username, password, loginChallenge string, httpHelper utils.HttpHelper) (string, error) {
	url := idmUrl + "/user/login"
	data := map[string]string{
		"username":  username,
		"password":  password,
		"challenge": loginChallenge,
	}
	json_data, err := json.Marshal(data)
	if err != nil {
		return "", errors.New("bad input body")
	}
	requestBody := bytes.NewBuffer(json_data)
	Requestheaders := []utils.Header{
		{
			Key:   "Content-Type",
			Value: "application/json",
		},
	}
	out, err, status, _ := httpHelper.HttpRequest(http.MethodPost, url, requestBody, Requestheaders, time.Minute, false)
	if err != nil {
		return "", err
	}
	if status != http.StatusOK {
		return "", errors.New("not Ok in post request, status: " + strconv.Itoa(status))
	}
	return string(out), nil
}

func getLocation(url, cookie string, httpHelper utils.HttpHelper) (string, []string, error) {
	Requestheaders := []utils.Header{
		{
			Key:   "Cookie",
			Value: cookie,
		},
	}
	_, err, status, headers := httpHelper.HttpRequest(http.MethodGet, url, nil, Requestheaders, time.Minute, false)
	if err != nil || status != http.StatusFound {
		var errString string
		if err != nil {
			errString = err.Error()
		}
		return "", nil, errors.New("error in getLocation request, err= " + errString + " status:" + strconv.Itoa(status))
	}
	if headers.Get("location") == "" {
		return "", nil, errors.New("no location header")
	}
	return headers.Get("location"), headers.Values("Set-Cookie"), nil
}
func getParam(key string, params []string) string {
	for _, param := range params {
		if strings.Contains(param, key) {
			return strings.Replace(param, key, "", 1)
		}
	}
	return ""
}
