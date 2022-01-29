/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"utopiops-cli/utils"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const ()

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list command shows list of applications, environments or ...",
	Long:  `list command shows list of applications, environments or ... based on what you pass with --all flag`,
	Run: func(cmd *cobra.Command, args []string) {
		what, err := cmd.Flags().GetString("all")
		if err != nil {
			fmt.Println("error in getting flags")
		}
		switch what {
		case "applications":
			err = showApps()
			if err != nil {
				if err.Error() == "not ok with status: 401" {
					RegisterCli()
					err = showApps()
					if err != nil {
						utils.HandleError(err, "showing environments list")
					}
				} else {
					utils.HandleError(err, "showing applications list")
				}
			}
		case "environments":
			err = showEnvs()
			if err != nil {
				if err.Error() == "not ok with status: 401" {
					RegisterCli()
					err = showEnvs()
					if err != nil {
						utils.HandleError(err, "showing environments list")
					}
				} else {
					utils.HandleError(err, "showing applications list")
				}
			}
		default:
			fmt.Println("we don't support showing you a list of " + what + " for now")
		}
	},
}

func showApps() error {
	blue := color.New(color.BgHiBlue).SprintFunc()
	green := color.New(color.BgGreen).SprintFunc()
	orrange := color.New(color.BgYellow).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	apps, err := UtopiopsService.GetApplications(viper.GetString("UTOPIOPS_OUTH_TOKEN"), viper.GetString("UTOPIOPS_OUTH_ID_TOKEN"))
	if err != nil {
		return err
	}
	for _, app := range apps {
		var status string
		switch app.Status {
		case "healthy":
			status = green(app.Status)
		case "warning":
			status = orrange(app.Status)
		case "no_alarm":
			status = blue(app.Status)
		case "critical":
			status = red(app.Status)
		}
		fmt.Printf("Name: %s\t kind: %s\t state: %s\t environment: %s\t status: %s\n", app.Name, app.Kind, app.State.Code, app.EnviironmentName, status)
	}
	return nil
}
func showEnvs() error {
	blue := color.New(color.BgHiBlue).SprintFunc()
	green := color.New(color.BgGreen).SprintFunc()
	orrange := color.New(color.BgYellow).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	apps, err := UtopiopsService.GetEnvironments(viper.GetString("UTOPIOPS_OUTH_TOKEN"), viper.GetString("UTOPIOPS_OUTH_ID_TOKEN"))
	if err != nil {
		return err
	}
	for _, app := range apps {
		var status string
		switch app.Status {
		case "healthy":
			status = green(app.Status)
		case "warning":
			status = orrange(app.Status)
		case "no_alarm":
			status = blue(app.Status)
		case "critical":
			status = red(app.Status)
		}
		fmt.Printf("Name: %s\t kind: %s\t state: %s\t provider: %s\t status: %s\n", app.Name, app.Kind, app.State.Code, app.ProviderName, status)
	}
	return nil
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringP("all", "a", "applications", "Name of what you want a list of it")
}
