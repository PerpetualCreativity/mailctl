package cmd

import (
	"os"
	_ "embed"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

//go:embed mailctl.yml
var defconfig []byte

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Creates a config file.",
	Long: `'mailctl init' creates a configuration file
	for your site.`,
	Run: func(cmd *cobra.Command, args []string) {
		homedir, err := homedir.Dir()
		fc.ErrCheck(err, "Could not find home directory")

		err = os.WriteFile(homedir+"/.mailctl.yml", defconfig, 0644)
		fc.ErrCheck(err, "Could not successfully create config file")

		fc.Success("Created configuration file at " + homedir+"/.mailctl.yml")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
