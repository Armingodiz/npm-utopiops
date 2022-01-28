/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"utopiops-cli/models"
	"utopiops-cli/utils"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "create static website",
	Long:  `command to create and deploy static website`,
	Run: func(cmd *cobra.Command, args []string) {
		//name, err1 := cmd.Flags().GetString("name")
		isCustom, err1 := cmd.Flags().GetBool("is-custom-domain")
		appType, err2 := cmd.Flags().GetString("type")
		if err1 != nil || err2 != nil {
			fmt.Println("error in getting flags")
		} else {
			if !checkType(appType) {
				fmt.Println("unsupported application type")
				return
			}
			in := bufio.NewReader(os.Stdin)
			fmt.Print("Enter name of application: ")
			name, _ := in.ReadString('\n')
			name = strings.Replace(name, "\n", "", 1)
			fmt.Print("Enter Description of application(can be empty): ")
			description, _ := in.ReadString('\n')
			description = strings.Replace(description, "\n", "", 1)
			switch appType {
			case "static-website":
				{
					cr := getStaticWebsiteCredentials(in)
					cr.Name = name
					cr.Description = description
					if isCustom {
						fmt.Print("Enter the domain: ")
						domain, _ := in.ReadString('\n')
						cr.Domain = strings.Replace(domain, "\n", "", 1)
						err := UtopiopsService.CreateCustomDomainStaticWebsite(cr, viper.GetString("UTOPIOPS_OUTH_TOKEN"), viper.GetString("UTOPIOPS_OUTH_ID_TOKEN"))
						if err != nil {
							if err.Error() == "not ok with status: 401" {
								RegisterCli()
							} else {
								utils.HandleError(err, "creating static website with custom domain")
								return
							}
							err = UtopiopsService.CreateCustomDomainStaticWebsite(cr, viper.GetString("UTOPIOPS_OUTH_TOKEN"), viper.GetString("UTOPIOPS_OUTH_ID_TOKEN"))
							if err != nil {
								utils.HandleError(err, "creating static website with custom domain")
								return
							}
						}
					} else {
						err := UtopiopsService.CreateStaticWebsite(cr, viper.GetString("UTOPIOPS_OUTH_TOKEN"), viper.GetString("UTOPIOPS_OUTH_ID_TOKEN"))
						if err != nil {
							if err.Error() == "not ok with status: 401" {
								RegisterCli()
							} else {
								utils.HandleError(err, "creating static website")
								return
							}
							err = UtopiopsService.CreateStaticWebsite(cr, viper.GetString("UTOPIOPS_OUTH_TOKEN"), viper.GetString("UTOPIOPS_OUTH_ID_TOKEN"))
							if err != nil {
								utils.HandleError(err, "creating static website")
								return
							}
						}
					}
					domain, err := UtopiopsService.GetStaticWebsiteDomain(name, viper.GetString("UTOPIOPS_OUTH_TOKEN"), viper.GetString("UTOPIOPS_OUTH_ID_TOKEN"))
					if err != nil {
						utils.HandleError(err, "getting created domain")
					}
					fmt.Printf("static website created with url: %s\n", domain)
				}
			case "s3-website":
				{
					cr := getS3StaticWebsiteCredentials(name, description, in)
					err := UtopiopsService.CreateS3StaticWebsite(cr, viper.GetString("UTOPIOPS_OUTH_TOKEN"), viper.GetString("UTOPIOPS_OUTH_ID_TOKEN"))
					if err != nil {
						if err.Error() == "not ok with status: 401" {
							RegisterCli()
						} else {
							utils.HandleError(err, "creating s3 static website")
							return
						}
						err = UtopiopsService.CreateS3StaticWebsite(cr, viper.GetString("UTOPIOPS_OUTH_TOKEN"), viper.GetString("UTOPIOPS_OUTH_ID_TOKEN"))
						if err != nil {
							utils.HandleError(err, "creating s3 static website")
							return
						}
					}
					fmt.Println("s3-website created successfully")
				}
			case "ecsapp":
				{
					cr := getEcsAppCredentials(name, description, in)
					err := UtopiopsService.CreateEcsApplication(cr, viper.GetString("UTOPIOPS_OUTH_TOKEN"), viper.GetString("UTOPIOPS_OUTH_ID_TOKEN"))
					if err != nil {
						if err.Error() == "not ok with status: 401" {
							RegisterCli()
						} else {
							utils.HandleError(err, "creating ecs application")
							return
						}
						err = UtopiopsService.CreateEcsApplication(cr, viper.GetString("UTOPIOPS_OUTH_TOKEN"), viper.GetString("UTOPIOPS_OUTH_ID_TOKEN"))
						if err != nil {
							utils.HandleError(err, "creating ecs application")
							return
						}
					}
					fmt.Println("ecs app created successfully")
				}
			case "dockerized":
				{
					cr := getDockerizedCredentils(name, description, isCustom, in)
					err := UtopiopsService.CreateDockerized(cr, viper.GetString("UTOPIOPS_OUTH_TOKEN"), viper.GetString("UTOPIOPS_OUTH_ID_TOKEN"))
					if err != nil {
						if err.Error() == "not ok with status: 401" {
							RegisterCli()
						} else {
							utils.HandleError(err, "creating dockerizes application")
							return
						}
						err = UtopiopsService.CreateDockerized(cr, viper.GetString("UTOPIOPS_OUTH_TOKEN"), viper.GetString("UTOPIOPS_OUTH_ID_TOKEN"))
						if err != nil {
							utils.HandleError(err, "creating dockerizes application")
							return
						}
					}
					fmt.Println("dockerizes app created successfully")

				}
			case "function":
				{
					cr := getFunctionCredentials(name, description, isCustom, in)
					err := UtopiopsService.CreateFunction(cr, viper.GetString("UTOPIOPS_OUTH_TOKEN"), viper.GetString("UTOPIOPS_OUTH_ID_TOKEN"))
					if err != nil {
						if err.Error() == "not ok with status: 401" {
							RegisterCli()
						} else {
							utils.HandleError(err, "creating function")
							return
						}
						err = UtopiopsService.CreateFunction(cr, viper.GetString("UTOPIOPS_OUTH_TOKEN"), viper.GetString("UTOPIOPS_OUTH_ID_TOKEN"))
						if err != nil {
							utils.HandleError(err, "creating function")
							return
						}
					}
					fmt.Println("function created successfully")
					domain, err := UtopiopsService.GetStaticWebsiteDomain(name, viper.GetString("UTOPIOPS_OUTH_TOKEN"), viper.GetString("UTOPIOPS_OUTH_ID_TOKEN"))
					if err != nil {
						utils.HandleError(err, "getting created domain")
					}
					fmt.Printf("your function domain will be: %s\n", domain)
				}
			}
		}
	},
}

func getStaticWebsiteCredentials(in *bufio.Reader) models.StaticWebsiteCredentials {
	cr := models.StaticWebsiteCredentials{}
	fmt.Print("Enter the repository url: ")
	repo, _ := in.ReadString('\n')
	cr.RepositoryUrl = strings.Replace(repo, "\n", "", 1)

	fmt.Print("Enter the build command: ")
	build, _ := in.ReadString('\n')
	cr.BuildCommand = strings.Replace(build, "\n", "", 1)

	fmt.Print("Enter the output path: ")
	out, _ := in.ReadString('\n')
	cr.OutputPath = strings.Replace(out, "\n", "", 1)

	fmt.Print("Enter the index document(Press enter to use Default==> index.html): ")
	index, _ := in.ReadString('\n')
	cr.Index_document = strings.Replace(index, "\n", "", 1)

	fmt.Print("Enter the error document(Press enter to use Default==> error.html): ")
	erroro, _ := in.ReadString('\n')
	cr.Error_document = strings.Replace(erroro, "\n", "", 1)

	fmt.Print("Enter the branch(Press enter to use Default==> main): ")
	branch, _ := in.ReadString('\n')
	cr.Branch = strings.Replace(branch, "\n", "", 1)

	return cr
}
func getS3StaticWebsiteCredentials(name, desc string, in *bufio.Reader) models.S3StaticWebsiteCredentials {
	cr := models.S3StaticWebsiteCredentials{}
	app := models.S3Application{}
	app.Description = desc
	app.Name = name
	environment := models.Environment{}
	fmt.Print("Enter the repository url: ")
	repo, _ := in.ReadString('\n')
	app.RepositoryUrl = strings.Replace(repo, "\n", "", 1)

	fmt.Print("Enter the environment name: ")
	env, _ := in.ReadString('\n')
	environment.Name = strings.Replace(env, "\n", "", 1)
	app.Type = "s3-website"

	apps := make([]models.S3Application, 0)
	apps = append(apps, app)
	envs := make([]models.Environment, 0)
	envs = append(envs, environment)
	cr.Applications = apps
	cr.Environments = envs
	return cr
}
func getEcsAppCredentials(name, desc string, in *bufio.Reader) models.EcsApplicationCredentials {
	cr := models.EcsApplicationCredentials{}
	app := models.EcsApplication{}
	app.Description = desc
	app.Name = name
	environment := models.Environment{}
	fmt.Print("Enter the port: ")
	port, _ := in.ReadString('\n')
	app.Port = strings.Replace(port, "\n", "", 1)

	fmt.Print("Enter the protocol(Press enter to use Default==> https): ")
	protocol, _ := in.ReadString('\n')
	app.Protocol = strings.Replace(protocol, "\n", "", 1)

	fmt.Print("Enter the environment name: ")
	env, _ := in.ReadString('\n')
	environment.Name = strings.Replace(env, "\n", "", 1)
	app.Type = "ecs"

	apps := make([]models.EcsApplication, 0)
	apps = append(apps, app)
	envs := make([]models.Environment, 0)
	envs = append(envs, environment)
	cr.Applications = apps
	cr.Environments = envs
	return cr
}

func getDockerizedCredentils(name, desc string, isCustom bool, in *bufio.Reader) models.DockerizedCredentials {
	cr := models.DockerizedCredentials{}
	cr.Name = name
	cr.Description = desc
	fmt.Print("Enter the repository url: ")
	repo, _ := in.ReadString('\n')
	cr.Repository = strings.Replace(repo, "\n", "", 1)

	fmt.Print("Enter the branch: ")
	branch, _ := in.ReadString('\n')
	cr.Branch = strings.Replace(branch, "\n", "", 1)

	fmt.Print("Enter the port: ")
	port, _ := in.ReadString('\n')
	cr.Port = strings.Replace(port, "\n", "", 1)
	if isCustom {
		fmt.Print("Enter the domain name: ")
		dom, _ := in.ReadString('\n')
		cr.DomainName = strings.Replace(dom, "\n", "", 1)
		fmt.Print("Enter the cpu: ")
		fmt.Scanln(&cr.Cpu)
		fmt.Print("Enter the memory: ")
		fmt.Scanln(&cr.Memory)
	}
	return cr
}
func getFunctionCredentials(name, desc string, isCustom bool, in *bufio.Reader) models.FunctionCredentials {
	cr := models.FunctionCredentials{}
	cr.Name = name
	cr.Description = desc
	fmt.Print("Enter the repository url: ")
	repo, _ := in.ReadString('\n')
	cr.Repository = strings.Replace(repo, "\n", "", 1)

	fmt.Print("Enter the branch: ")
	branch, _ := in.ReadString('\n')
	cr.Branch = strings.Replace(branch, "\n", "", 1)
	if isCustom {
		fmt.Print("Enter the domain name: ")
		dom, _ := in.ReadString('\n')
		cr.DomainName = strings.Replace(dom, "\n", "", 1)
	}
	return cr
}

func checkType(t string) bool {
	if t == "ecsapp" || t == "static-website" || t == "s3-website" || t == "dockerized" || t == "function" {
		return true
	}
	return false
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().StringP("type", "t", "static-website", "type of application(static-website, ecsapp, s3-website, dockerized, function)")
	createCmd.Flags().Bool("is-custom-domain", false, "is domain custom(only used for static-website, dockerized and function types)")
}
