/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"utopiops-cli/models"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "deploy ecs",
	Long:  `command to deploy ecs apps on utopiops`,
	Run: func(cmd *cobra.Command, args []string) {
		app, err := cmd.Flags().GetString("application")
		environment, err2 := cmd.Flags().GetString("environment")
		tag, err3 := cmd.Flags().GetString("image-tag")
		repository, err4 := cmd.Flags().GetString("repo")
		if err != nil || err2 != nil || err3 != nil || err4 != nil {
			fmt.Println("error in getting flags")
		} else {
			err = deploy(app, environment, tag, repository)
			if err != nil {
				//fmt.Println(err)
				fmt.Println("error in deploying application")
				return
			}
		}
		fmt.Println("deploying started")
	},
}

func deploy(app, environment, tag_id, repo string) (err error) {
	detail, err := UtopiopsService.GetApplicationDetailes(viper.GetString("CORE_URL"), app,
		environment, viper.GetString("UTOPIOPS_OUTH_TOKEN"), viper.GetString("UTOPIOPS_OUTH_ID_TOKEN"))
	if err != nil {
		if err.Error() == "not ok with status: 401" {
			RegisterCli()
		} else {
			return errors.New("error in getting app details")
		}
		detail, err = UtopiopsService.GetApplicationDetailes(viper.GetString("CORE_URL"), app,
			environment, viper.GetString("UTOPIOPS_OUTH_TOKEN"), viper.GetString("UTOPIOPS_OUTH_ID_TOKEN"))
		if err != nil {
			return errors.New("error in getting app details")
		}
	}
	pushCredentials := models.PushCredentials{
		ImageTag:   tag_id,
		EcrUrl:     detail.EcrRegisteryUrl,
		Repository: repo,
	}
	err = pushImage(pushCredentials)
	if err != nil {
		fmt.Print("push step error")
		return
	}
	containerTags := make([]models.ContainerTag, 0)
	containerTags = append(containerTags, models.ContainerTag{
		ContainerName: detail.ContainerNames[0],
		ImageTag:      tag_id,
	})
	deployCredintials := models.DeployToUtopiopsCredentials{
		CoreUrl:      viper.GetString("CORE_URL"),
		Environment:  environment,
		Application:  app,
		Token:        viper.GetString("UTOPIOPS_OUTH_TOKEN"),
		IdToken:      viper.GetString("UTOPIOPS_OUTH_ID_TOKEN"),
		ContainerTag: containerTags,
	}
	return UtopiopsService.Deploy(deployCredintials)
}

func pushImage(cr models.PushCredentials) error {
	if err := cr.IsValid(); err != nil {
		return err
	}
	image := cr.EcrUrl + "/" + cr.Repository + ":" + cr.ImageTag
	cmd := exec.Command("docker", "push", image)
	out, err := cmd.CombinedOutput()
	if err != nil {
		if strings.Contains(string(out), "image does not exist") {
			return errors.New("image with tag [" + cr.ImageTag + "] does not exist")
		}
		return err
	}
	return nil
}

func init() {
	rootCmd.AddCommand(deployCmd)
	deployCmd.Flags().String("repo", "", "repositry name")
	deployCmd.Flags().String("image-tag", "", "image tag id")
	deployCmd.Flags().StringP("environment", "e", "", "environment name")
	deployCmd.Flags().StringP("application", "a", "", "Name of application to be deployed")
}
