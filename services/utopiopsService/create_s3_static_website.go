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
	"utopiops-cli/models"
	"utopiops-cli/utils"

	"github.com/r3labs/sse"
	"github.com/spf13/viper"
)

func (manager *UtopiopsManager) CreateS3StaticWebsite(cr models.S3StaticWebsiteCredentials, token, idToken string) error {
	cr = cr.SetDefaults()
	url := viper.GetString("DM_URL") + "/flash-setup/stage-2"
	return createWithLog(cr, url, token, idToken, manager.HttpHelper)
}

func createWithLog(cr models.CreateCredentials, url, token, idToken string, httpHelper utils.HttpHelper) error {
	if err := cr.IsValid(); err != nil {
		if err != nil {
			return err
		}
	}
	json_data, err := json.Marshal(cr)
	//dt, _ := json.MarshalIndent(cr, "", "    ")
	if err != nil {
		return errors.New("bad input body")
	}
	//fmt.Println(string(dt))
	requestBody := bytes.NewBuffer(json_data)
	Requestheaders := []utils.Header{
		{
			Key:   "Authorization",
			Value: fmt.Sprintf("Bearer %s", token),
		},
		{Key: "Cookie",
			Value: fmt.Sprintf("id_token=%s", idToken),
		},
		{
			Key:   "Content-Type",
			Value: "application/json",
		},
	}
	out, err, status, _ := httpHelper.HttpRequest(http.MethodPost, url, requestBody, Requestheaders, time.Minute, true)
	if err != nil {
		return err
	}
	//fmt.Println(string(out))
	if status != http.StatusOK && status != http.StatusAccepted {
		return errors.New("not ok with status: " + strconv.Itoa(status))
	}
	var res struct {
		JobId string `json:"jobId"`
	}
	if err := json.Unmarshal(out, &res); err != nil {
		return err
	}
	return getLogs(res.JobId, viper.GetString("LSM_URL"), token, idToken, httpHelper)
}
func getLogs(jobId, lsmUrl, token, idToken string, httpHelper utils.HttpHelper) error {
	url := fmt.Sprintf("%s/log/job?jobId=%s", lsmUrl, jobId)
	client := sse.NewClient(url)
	client.Headers = map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", token),
		"Cookie":        fmt.Sprintf("id_token=%s", idToken),
	}
	events := make(chan *sse.Event)
	client.SubscribeChan("messages", events)
	for {
		var current []Log
		event := <-events
		if string(event.Event) == "end" {
			return nil
		}
		if err := json.Unmarshal(event.Data, &current); err != nil {
			return err
		}
		err := processLog(current)
		if err != nil {
			return err
		}
	}
}

type Log struct {
	Id         string `json:"jobId"`
	Line       int    `json:"lineNumber"`
	Payload    string `json:"payload"`
	IsLastLine bool   `json:"isLastLine"`
}

func processLog(logs []Log) error {
	for _, log := range logs {
		fmt.Print(log.Line)
		fmt.Print(") ")
		fmt.Println(log.Payload)
		if strings.Contains(log.Payload, "an error occurred") {
			return errors.New("failed with error")
		}
	}
	return nil
}
