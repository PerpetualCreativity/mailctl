package cmd

import (
	"github.com/PerpetualCreativity/mailctl/utils"
	"strconv"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

// moveCmd represents the move command
var moveCmd = &cobra.Command{
	Use:   "move id currentFolder [toFolder]",
	Short: "move email or draft to another folder",
	Long: `Move email or draft to another folder. Requires id
and current folder. If toFolder is not specified, options
will be listed.`,
	Args: cobra.RangeArgs(2, 3),
	Run: func(cmd *cobra.Command, args []string) {
		c, err := utils.ImapLogin()
		fc.ErrCheck(err, "error when logging in to IMAP server")
		defer c.Logout()

		id, err := strconv.Atoi(args[0])
		fc.ErrCheck(err, "ID is not an integer")

		// prompt user for folder names if necessary
		mailboxNames := utils.ListMailboxes(c)

		toFolder := ""
		if len(args) <= 2 {
			folderPrompt := &survey.Select{
				Message:  "Move to folder: ",
				Options:  mailboxNames,
				Default:  utils.FindMailbox(c, "\\Inbox", "Inbox"),
				PageSize: 10,
			}
			err := survey.AskOne(folderPrompt, &toFolder)
			fc.ErrCheck(err, "Prompt failed")
		} else {
			toFolder = args[2]
		}

		err = utils.MoveMail(c, id, args[1], toFolder)
		fc.ErrCheck(err, "error when moving mail")

		fc.Success("Moved to " + toFolder)
	},
}

func init() {
	rootCmd.AddCommand(moveCmd)
}
