package cmd

import (
	"github.com/PerpetualCreativity/mailctl/utils"
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"

	"github.com/AlecAivazis/survey/v2"
	"github.com/emersion/go-imap"
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
		c := utils.ImapLogin()
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
			utils.ErrCheck(err, "Prompt failed")
		} else {
			folder = args[0]
		}

		// select mailbox
		mailbox, err := c.Select(folder, false)
		utils.ErrCheck(err, "Could not select mailbox")
		// get last numberMessages messages
		from := uint32(1)
		to := mailbox.Messages

		if mailbox.Messages > numberMessages {
			from = mailbox.Messages - numberMessages
		}
		seqset := new(imap.SeqSet)
		seqset.AddRange(from, to)

		messages := make(chan *imap.Message, numberMessages)
		done := make(chan error, 1)
		go func() {
			done <- c.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope}, messages)
		}()

		fmt.Printf("Last %d messages:\n", numberMessages)
		tw := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
		for msg := range messages {
			sender := ""
			if s := msg.Envelope.From; len(s) > 0 {
				sender = msg.Envelope.From[0].PersonalName
			}
			fmt.Fprintf(tw, "#%s\t%s\t %s\n", strconv.FormatUint(uint64(msg.SeqNum), 10), sender, msg.Envelope.Subject)
		}
		tw.Flush()

		utils.ErrCheck(<-done, "No messsages in this folder")
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
