package cmd

import (
	"bytes"
	"github.com/PerpetualCreativity/mailctl/utils"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Write a new email.",
	Long:  `Write a new email. Call the command to try it out...`,
	Run: func(cmd *cobra.Command, args []string) {
		// log in to IMAP server
		c := utils.ImapLogin()
		defer c.Logout()

		messagePrompt := []*survey.Question{
			{
				Name: "body",
				Prompt: &survey.Editor{
					Message:  "Content",
					FileName: "*.md",
				},
				Validate: survey.Required,
			},
			{
				Name:     "subject",
				Prompt:   &survey.Input{Message: "Subject: "},
				Validate: survey.Required,
			},
		}

		response := struct {
			Body    string
			Subject string
		}{}

		survey.Ask(messagePrompt, &response)

		var msg bytes.Buffer
		msg.WriteString("Subject: " + response.Subject + "\r\n")
		msg.WriteString("\r\n")
		msg.WriteString(response.Body)

		// find which mailbox is for drafts
		draftsBox := utils.FindMailbox(c, "\\Drafts", "Drafts")
		err := c.Append(draftsBox, nil, time.Now(), &msg)
		fc.ErrCheck(err, "Could not add draft to Drafts folder")

		fc.Success("Created draft.")
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
}
