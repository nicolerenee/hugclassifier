package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var (
	cfgFile   string
	opAccount string
	opLogin   string
	opVault   string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "hugclassifier",
	Short: "Return HUG meetup data from meetup API",
	Long: `HUG Classifier is an application that talks to the Meetup API
and pulls data about meetups happening within the HUG Pro Network. Some
examples of how to use the application are as follows:

To get the last 30 days (default) of meetup events:
	hugclassifier generate

To get the last 90 days of meetup events:
	hugclassifier generate --days 90

To get a specific date range of meetup events:
	hugclassifier generate --start 2020-01-01 --end 2020-04-01

`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.hugclassifier.yaml)")

	rootCmd.PersistentFlags().String("onepass-account", "hashicorp", "1Password account name")
	viper.BindPFlag("onepassword.account", rootCmd.PersistentFlags().Lookup("onepass-account"))
	rootCmd.PersistentFlags().String("onepass-vault", "", "1Password vault UUID with meetup API credentials")
	viper.BindPFlag("onepassword.vault", rootCmd.PersistentFlags().Lookup("onepass-vault"))
	rootCmd.PersistentFlags().String("onepass-login", "", "1Password login UUID with meetup API credentials")
	viper.BindPFlag("onepassword.login", rootCmd.PersistentFlags().Lookup("onepass-login"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".hugclassifier" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".hugclassifier")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Printf("info: using config file: %s\n", viper.ConfigFileUsed())
	}
}
