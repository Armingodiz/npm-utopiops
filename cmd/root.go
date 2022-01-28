/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"utopiops-cli/services/awsService"
	"utopiops-cli/services/utopiopsService"
	"utopiops-cli/utils"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

// rootCmd represents the base command when called without any subcommands
var (
	UtopiopsService utopiopsService.UtopiopsService
	AwsService      awsService.AwsService
	cfgFile         string
	rootCmd         = &cobra.Command{
		Use:   "cli",
		Short: "Uopiops Cli",
		Long:  `cli tool for utopiops`,
		// Uncomment the following line if your bare application
		// has an action associated with it:
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Utopiops is here!")
		},
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	cobra.OnInitialize(initConfig)

	//RegisterCli()
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.utopiops.yaml) will be created for you if it doesn't exists")
	//rootCmd.PersistentFlags().String("AWS_ACCESS_KEY_ID", "", "aws access key id")
	//rootCmd.PersistentFlags().String("AWS_SECRET_ACCESS_KEY", "", "aws secret access key")
	//rootCmd.PersistentFlags().String("UTOPIOPS_OUTH_TOKEN", "", "token for intracting with utopiops apis")
	//rootCmd.PersistentFlags().String("UTOPIOPS_OUTH_ID_TOKEN", "", "id token for intracting with utopiops apis")
	//rootCmd.PersistentFlags().String("CORE_URL", "", "url for core")
	//rootCmd.PersistentFlags().StringP("application", "a", "", "Name of application to be build/deployed")
	//rootCmd.PersistentFlags().String("IDS_URL", "", "ids url")
	//rootCmd.PersistentFlags().String("IDM_URL", "", "idm url")
	//rootCmd.PersistentFlags().Bool("viper", true, "use Viper for configuration")
	// bind flags with viper
	//viper.BindPFlag("AWS_ACCESS_KEY_ID", rootCmd.PersistentFlags().Lookup("AWS_ACCESS_KEY_ID"))
	//viper.BindPFlag("AWS_SECRET_ACCESS_KEY", rootCmd.PersistentFlags().Lookup("AWS_SECRET_ACCESS_KEY"))
	//viper.BindPFlag("UTOPIOPS_OUTH_TOKEN", rootCmd.PersistentFlags().Lookup("UTOPIOPS_OUTH_TOKEN"))
	//viper.BindPFlag("UTOPIOPS_OUTH_ID_TOKEN", rootCmd.PersistentFlags().Lookup("UTOPIOPS_OUTH_ID_TOKEN"))
	//viper.BindPFlag("CORE_URL", rootCmd.PersistentFlags().Lookup("CORE_URL"))
	//viper.BindPFlag("IDS_URL", rootCmd.PersistentFlags().Lookup("IDS_URL"))
	//viper.BindPFlag("IDM_URL", rootCmd.PersistentFlags().Lookup("IDM_URL"))
	//viper.BindPFlag("useViper", rootCmd.PersistentFlags().Lookup("viper"))
	httpHelper := utils.NewHttpHelper(utils.NewHttpClient())
	AwsService = awsService.NewService()
	UtopiopsService = utopiopsService.NewService(httpHelper, AwsService)
}

func initConfig() {
	//fmt.Println("in init config")
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".utopiops" (without extension).
		// todo remove /Desktop which is for test
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".utopiops")
		if runtime.GOOS == "windows" {
			cfgFile = home + "\\.utopiops.yml"
		} else {
			cfgFile = home + "/.utopiops.yml"
		}
		//log.Println("file set to be: " + cfgFile)
	}
	_, err := os.Open(cfgFile)
	if os.IsNotExist(err) {
		err = addConfigFile(cfgFile)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		fmt.Println(err)
	}
	if viper.GetString("UTOPIOPS_OUTH_TOKEN") == "" || viper.GetString("UTOPIOPS_OUTH_ID_TOKEN") == "" {
		err := RegisterCli()
		if err != nil {
			fmt.Println("registeration failed" + err.Error())
			os.Exit(1)
		}
	}
}

func RegisterCli() error {
	fmt.Println("registering cli ...")
	user, pass := viper.GetString("UTOPIOPS_USERNAME"), viper.GetString("UTOPIOPS_PASSWORD")
	if user == "" || pass == "" {
		fmt.Print("Enter your username: ")
		fmt.Scanln(&user)
		fmt.Print("Enter your password: ")
		fmt.Scanln(&pass)
	}
	token, id, err := UtopiopsService.Register(
		viper.GetString("IDS_URL"), viper.GetString("IDM_URL"),
		user, pass)
	if err == nil {
		var config Config
		yamlFile, err := ioutil.ReadFile(cfgFile)
		if err != nil {
			log.Printf("yamlFile.Get err   #%v ", err)
		}
		err = yaml.Unmarshal(yamlFile, &config)
		if err != nil {
			log.Fatalf("Unmarshal: %v", err)
		}
		config.UtopiopsOuthToken = token
		viper.Set("UTOPIOPS_OUTH_TOKEN", token)
		config.UtopiopsOuthIdToken = id
		viper.Set("UTOPIOPS_OUTH_ID_TOKEN", id)
		config.UtopiopsUserName = user
		config.UtopiopsPassword = pass
		marshalledConfig, _ := yaml.Marshal(&config)
		return ioutil.WriteFile(cfgFile, marshalledConfig, fs.ModeAppend)
	}
	return err
}

func addConfigFile(path string) error {
	var environment string
	fmt.Println("what environment is your workground?(press enter if you are not a utopiops develpoer)")
	fmt.Scanln(&environment)
	var config Config
	if environment == "staging" {
		config = Config{
			CoreUrl: "https://core.staging.utopiops.com",
			IdsUrl:  "https://ids-pub.staging.utopiops.com",
			IdmUrl:  "https://idm.staging.utopiops.com",
			LSMUrl:  "https://lsm.staging.utopiops.com",
			DMUrl:   "https://dm.staging.utopiops.com",
		}
	} else if environment == "production" || environment == "" {
		config = Config{
			CoreUrl: "https://core.utopiops.com",
			IdsUrl:  "https://ids.utopiops.com",
			IdmUrl:  "https://idm.utopiops.com",
			LSMUrl:  "https://lsm.utopiops.com",
			DMUrl:   "https://dm.utopiops.com",
		}
	} else {
		return fmt.Errorf("environment %s doesn't exists", environment)
	}
	marshalledConfig, _ := yaml.Marshal(&config)
	_, err := os.Create(path)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, marshalledConfig, fs.ModeAppend)
}
