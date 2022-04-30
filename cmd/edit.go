package cmd

import (
	"bytes"
	"strconv"
	"strings"
	"time"

	"github.com/PerpetualCreativity/mailctl/utils"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:   "edit id [folder]",
	Short: "Edit a draft.",
	Long:  `Edit a draft. Call the command to try it out...`,
	Args: cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		// log in to IMAP server
		var id int
		id, err := strconv.Atoi(args[0])
		fc.ErrCheck(err, "ID is not an integer")

		c := utils.ImapLogin()
		defer c.Logout()

		folder := "Drafts"
		if len(args) > 1 {
			folder = args[1]
		}

		subject, body := utils.GetMessage(c, id, folder)

		messagePrompt := []*survey.Question{
			{
				Name: "body",
				Prompt: &survey.Editor{
					Message:  "Content",
					FileName: "*.md",
					Default: body,
					AppendDefault: true,
					HideDefault: true,
				},
				Validate: survey.Required,
			},
			{
				Name:     "subject",
				Prompt:   &survey.Input{
					Message: "Subject: ",
					Default: subject,
				},
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
		err = c.Append(draftsBox, nil, time.Now(), &msg)
		fc.ErrCheck(err, "Could not add draft to Drafts folder")
		utils.MoveMail(c, id, folder, "Deleted")

		fc.Success("Created draft.")
	},
}

func init() {
	rootCmd.AddCommand(editCmd)
}
