package cmd

import (
	"fmt"
	"github.com/PerpetualCreativity/mailctl/utils"
	"os"
	"strconv"
	"text/tabwriter"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list [folder] [number]",
	Short: "List emails in folder",
	Long: `This lists emails in a folder. If you
don't give a folder name, you will be presented with
options. Number is the number of messages to display
(default 10)`,
	Args: cobra.MaximumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		c, err := utils.ImapLogin()
		fc.ErrCheck(err, "error when logging into IMAP server")
		defer c.Logout()

		numberMessages := uint32(10)
		if len(args) > 1 {
			s, _ := strconv.Atoi(args[1])
			numberMessages = uint32(s)
		}

		var folder string
		// prompt user for folder name if necessary
		if len(args) == 0 || args[0] == "" {
			mailboxNames := utils.ListMailboxes(c)

			folderPrompt := &survey.Select{
				Message:  "Select a folder:",
				Options:  mailboxNames,
				Default:  "inbox",
				PageSize: 10,
			}
			err := survey.AskOne(folderPrompt, &folder)
			fc.ErrCheck(err, "Prompt failed")
		} else {
			folder = args[0]
		}

		messages, err := utils.ListMessages(c, folder, numberMessages, 0)
		fc.ErrCheck(err, "error when getting message list")

		if len(messages) == 0 {
			fc.Neutral("No messages in this folder.")
			return
		}

		fc.Neutral("Last " + fmt.Sprintf("%d", numberMessages) + " messages:\n")
		tw := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
		for _, msg := range messages {
			fmt.Fprintf(tw, "#%s\t%s\t %s\n", strconv.FormatUint(uint64(msg.SeqNum), 10), msg.Sender, msg.Subject)
		}
		tw.Flush()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
