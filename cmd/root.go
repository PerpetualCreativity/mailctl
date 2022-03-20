package cmd

import (
	"github.com/fatih/color"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/PerpetualCreativity/fancyChecks"
)

var fc = fancyChecks.New("", "", "Status", "Error")

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mailctl",
	Short: "A modern console-based mail application.",
	Long:  `mailctl is an intuitive console-based mail application.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	var account int
	rootCmd.PersistentFlags().IntVar(&account, "account", -1, "account to use for command (default is 1, if not set in config)")

	if account != -1 {
		viper.Set("default_account", account)
	} else if viper.IsSet("default_account") {
		viper.Set("default_account", viper.GetInt("default_account"))
	} else {
		viper.Set("default_account", 1)
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// Find home directory.
	home, err := homedir.Dir()
	cobra.CheckErr(err)

	viper.AddConfigPath(home)
	viper.SetConfigType("yaml")
	viper.SetConfigName(".mailctl")


	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		color.Blue("Could not find config file: ", err)
	}
}
