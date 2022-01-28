/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"utopiops-cli/models"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// logCmd represents the log command

var logCmd = &cobra.Command{
	Use:   "log",
	Short: "show logs",
	Long:  `log command show live logs of ecs application using cloudwatch`,
	Run: func(cmd *cobra.Command, args []string) {
		app, _ := cmd.Flags().GetString("application")
		env, _ := cmd.Flags().GetString("environment")
		exept, _ := cmd.Flags().GetString("exept")
		find, _ := cmd.Flags().GetString("find")
		from, _ := cmd.Flags().GetInt("from")
		l := models.Log{
			App:         app,
			Environment: env,
			Exept:       exept,
			Find:        find,
			From:        from,
		}
		if err := l.IsValid(); err != nil {
			fmt.Println(err.Error())
		} else {
			var err error
			if from > 0 {
				err = UtopiopsService.Show(l, viper.GetString("UTOPIOPS_OUTH_TOKEN"), viper.GetString("UTOPIOPS_OUTH_ID_TOKEN"))
			} else {
				err = UtopiopsService.Watch(l, viper.GetString("UTOPIOPS_OUTH_TOKEN"), viper.GetString("UTOPIOPS_OUTH_ID_TOKEN"))
			}
			if err != nil {
				if err.Error() == "not ok with status: 401" {
					RegisterCli()
					if from > 0 {
						err = UtopiopsService.Show(l, viper.GetString("UTOPIOPS_OUTH_TOKEN"), viper.GetString("UTOPIOPS_OUTH_ID_TOKEN"))
					} else {
						err = UtopiopsService.Watch(l, viper.GetString("UTOPIOPS_OUTH_TOKEN"), viper.GetString("UTOPIOPS_OUTH_ID_TOKEN"))
					}
				}
			}
			if err != nil {
				if err.Error() == "not ok with status: 401" {
					fmt.Println("error in registering cli, try again")
				}
				fmt.Println(err.Error())
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(logCmd)
	logCmd.Flags().StringP("application", "a", "", "Name of application to be deployed")
	logCmd.Flags().StringP("environment", "e", "", "environment name")
	logCmd.Flags().String("exept", "", "logs which contain this word won't be printed")
	logCmd.Flags().String("find", "", "logs which contain this word will be printed")
	logCmd.Flags().Int("from", 0, "if you set a number for from cli shows you logs from that number * minute ago")
}
