package awsService

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
	"utopiops-cli/models"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/fatih/color"
)

func (manager *AwsManager) Show(lrc models.Log, cr models.ProviderCredentials) error {
	if lrc.From == 0 {
		return errors.New("from is not valid")
	}
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{Region: &cr.Region, Credentials: credentials.NewStaticCredentials(cr.AccessKeyId, cr.SecretAccessKey, string(""))},
	}))
	svc := cloudwatchlogs.New(sess)

	// This context controls the overall execution of getting logs and sending the events. We intentionally stop the log stream after a duration.
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(80)*time.Second)
	return getLogs(lrc, ctx, &lrc.LogGroup, svc)
}

func getLogs(lrc models.Log, ctx context.Context, logGroupName *string, svc *cloudwatchlogs.CloudWatchLogs) (err error) {
	red := color.New(color.FgRed).SprintFunc()
	for {
		perQueryCtx, cancel := context.WithTimeout(context.Background(), time.Duration(40)*time.Second)
		input := cloudwatchlogs.StartQueryInput{
			LogGroupName: logGroupName,
			Limit:        aws.Int64(10),
			StartTime:    aws.Int64(time.Now().Add(time.Duration(-1*lrc.From) * time.Minute).Unix()),
			EndTime:      aws.Int64(time.Now().Unix()),
			QueryString:  aws.String("fields toMillis(@timestamp) as timestamp, @message"),
		}
		if err = input.Validate(); err != nil {
			cancel()
			return
		}

		req, startQueryOutput := svc.StartQueryRequest(&input)
		err = req.Send()
		if err != nil {
			cancel()
			return
		}
		queryResults := make(chan *cloudwatchlogs.GetQueryResultsOutput)
		queryErrors := make(chan error)
		go getQueryResultsUntilComplete(perQueryCtx, svc, *startQueryOutput.QueryId, 10, queryResults, queryErrors)
		select {
		case <-ctx.Done():
			cancel()
			break
		case <-perQueryCtx.Done():
			err = errors.New("query timed out")
			break
		case err = <-queryErrors:
			break
		case queryResult := <-queryResults:
			var dateStr, message string
			for _, rs := range queryResult.Results {
				for _, msgs := range rs {
					if *msgs.Field == "timestamp" {
						t, err := strconv.Atoi(*msgs.Value)
						if err == nil {
							tr := time.Unix(int64(t), 0)
							dateStr = fmt.Sprintf("[%s] ", red(tr.UTC()))
						} else {
							dateStr = fmt.Sprintf("[%s] ", red("Invalid timestamp"))
						}
					}
					if *msgs.Field == "@message" {
						jl := map[string]interface{}{}
						if err := json.Unmarshal([]byte(*msgs.Value), &jl); err != nil || len(jl) == 0 {
							message = *msgs.Value
						} else {
							dt, _ := json.MarshalIndent(jl, "", "    ")
							message = string(dt)
						}
						if lrc.Exept != "" && !strings.Contains(message, lrc.Exept) {
							fmt.Println(dateStr + message)
						} else if lrc.Find != "" && strings.Contains(message, lrc.Find) {
							fmt.Println(dateStr + message)
						} else if lrc.Find == "" && lrc.Exept == "" {
							fmt.Println(dateStr + message)
						}
					}
				}
			}
			break
		}
	}
}

func getQueryResultsUntilComplete(ctx context.Context, cwl *cloudwatchlogs.CloudWatchLogs, queryId string, limit int, results chan *cloudwatchlogs.GetQueryResultsOutput, errorsChan chan error) {
	getQueryResultInput := &cloudwatchlogs.GetQueryResultsInput{}
	getQueryResultInput.SetQueryId(queryId)
	for {
		getQueryResultOutput, err := cwl.GetQueryResultsWithContext(ctx, getQueryResultInput)
		if err != nil {
			errorsChan <- err
			return
		}
		time.Sleep(5 * time.Second)
		switch *getQueryResultOutput.Status {
		case "Running":
			if len(getQueryResultOutput.Results) < limit {
				continue
			}
			stopQueryInput := &cloudwatchlogs.StopQueryInput{}
			stopQueryInput.SetQueryId(queryId)
			stopResult, err := cwl.StopQuery(stopQueryInput)
			if err != nil {
				errorsChan <- fmt.Errorf("stop query error=%s status=%v", err.Error(), stopResult)
				return
			}
			results <- getQueryResultOutput
			close(results)
			return
		case "Scheduled":
			continue
		case "Failed":
			errorsChan <- errors.New("job failed")
			return
		case "Cancelled":
			errorsChan <- errors.New("job cancelled")
			return
		case "Complete":
			results <- getQueryResultOutput
			close(results)
			return
		default:
			errorsChan <- fmt.Errorf("unknown status: %s", *getQueryResultOutput.Status)
			return
		}
	}
}
