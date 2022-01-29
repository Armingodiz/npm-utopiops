/*
Copyright (c) 2017 Tyler Brock
This part of code is almost copied from https://github.com/TylerBrock/saw
Writer: https://github.com/TylerBrock
*/

package awsService

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
	"utopiops-cli/models"

	"github.com/TylerBrock/colorjson"
	"github.com/TylerBrock/saw/config"
	"github.com/fatih/color"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
)

var watchConfig config.Configuration

var watchOutputConfig config.OutputConfiguration

func (manager *AwsManager) Watch(lrc models.Log, cr models.ProviderCredentials) error {
	watchConfig.Group = lrc.LogGroup
	b := NewBlade(&watchConfig, &watchOutputConfig, aws.Config{Region: &cr.Region, Credentials: credentials.NewStaticCredentials(cr.AccessKeyId, cr.SecretAccessKey, string(""))})
	b.StreamEvents(lrc.Exept, lrc.Find)
	return nil
}

// A Blade is a Saw execution instance
type Blade struct {
	config *config.Configuration
	output *config.OutputConfiguration
	cwl    *cloudwatchlogs.CloudWatchLogs
}

// NewBlade creates a new Blade with CloudWatchLogs instance from provided config
func NewBlade(
	config *config.Configuration,
	outputConfig *config.OutputConfiguration,
	awsConfig aws.Config) *Blade {
	blade := Blade{}

	awsSessionOpts := session.Options{
		Config:            awsConfig,
		SharedConfigState: session.SharedConfigEnable,
	}

	sess := session.Must(session.NewSessionWithOptions(awsSessionOpts))
	blade.cwl = cloudwatchlogs.New(sess)
	blade.config = config
	blade.output = outputConfig

	return &blade
}

// GetLogStreams gets the log streams from AWS given the blade configuration
func (b *Blade) GetLogStreams() []*cloudwatchlogs.LogStream {
	input := b.config.DescribeLogStreamsInput()
	streams := make([]*cloudwatchlogs.LogStream, 0)
	b.cwl.DescribeLogStreamsPages(input, func(
		out *cloudwatchlogs.DescribeLogStreamsOutput,
		lastPage bool,
	) bool {
		streams = append(streams, out.LogStreams...)
		return !lastPage
	})
	return streams
}

// GetEvents gets events from AWS given the blade configuration
func (b *Blade) GetEvents() {
	formatter := b.output.Formatter()
	input := b.config.FilterLogEventsInput()

	handlePage := func(page *cloudwatchlogs.FilterLogEventsOutput, lastPage bool) bool {
		for _, event := range page.Events {
			if b.output.Pretty {
				fmt.Println(formatEvent(formatter, event))
			} else {
				fmt.Println(*event.Message)
			}
		}
		return !lastPage
	}
	err := b.cwl.FilterLogEventsPages(input, handlePage)
	if err != nil {
		fmt.Println("Error", err)
		os.Exit(2)
	}
}

// StreamEvents continuously prints log events to the console
func (b *Blade) StreamEvents(exept, find string) {
	var lastSeenTime *int64
	var seenEventIDs map[string]bool
	formatter := b.output.Formatter()
	input := b.config.FilterLogEventsInput()

	clearSeenEventIds := func() {
		seenEventIDs = make(map[string]bool, 0)
	}

	addSeenEventIDs := func(id *string) {
		seenEventIDs[*id] = true
	}

	updateLastSeenTime := func(ts *int64) {
		if lastSeenTime == nil || *ts > *lastSeenTime {
			lastSeenTime = ts
			clearSeenEventIds()
		}
	}

	handlePage := func(page *cloudwatchlogs.FilterLogEventsOutput, lastPage bool) bool {
		for _, event := range page.Events {
			updateLastSeenTime(event.Timestamp)
			if _, seen := seenEventIDs[*event.EventId]; !seen {
				var message string
				if b.output.Raw {
					message = *event.Message
				} else {
					message = formatEvent(formatter, event)
				}
				message = strings.TrimRight(message, "\n")
				if exept != "" && !strings.Contains(message, exept) {
					fmt.Println(message)
				} else if find != "" && strings.Contains(message, find) {
					fmt.Println(message)
				} else if find == "" && exept == "" {
					fmt.Println(message)
				}
				addSeenEventIDs(event.EventId)
			}
		}
		return !lastPage
	}

	for {
		err := b.cwl.FilterLogEventsPages(input, handlePage)
		if err != nil {
			fmt.Println("Error", err)
			os.Exit(2)
		}
		if lastSeenTime != nil {
			input.SetStartTime(*lastSeenTime)
		}
		time.Sleep(1 * time.Second)
	}
}

// formatEvent returns a CloudWatch log event as a formatted string using the provided formatter
func formatEvent(formatter *colorjson.Formatter, event *cloudwatchlogs.FilteredLogEvent) string {
	red := color.New(color.FgRed).SprintFunc()
	white := color.New(color.FgWhite).SprintFunc()

	str := aws.StringValue(event.Message)
	bytes := []byte(str)
	date := aws.MillisecondsTimeValue(event.Timestamp)
	dateStr := date.Format(time.RFC3339)
	//streamStr := aws.StringValue(event.LogStreamName)
	jl := map[string]interface{}{}

	if err := json.Unmarshal(bytes, &jl); err != nil {
		return fmt.Sprintf("[%s] %s ", red(dateStr), white(str))
	}

	//output, _ := formatter.Marshal(jl)
	dt, err := json.MarshalIndent(jl, "", "    ")
	if err != nil {
		return fmt.Sprintf("[%s] %s", red(dateStr), "error in marshalling message, it is not in json format, error: "+err.Error())
	}
	return fmt.Sprintf("[%s] %s", red(dateStr), dt)
}
