package cmd

import (
	"github.com/PerpetualCreativity/mailctl/utils"
	"strconv"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

// sendCmd represents the send command
var sendCmd = &cobra.Command{
	Use:   "send id",
	Short: "Send an email (requires draft id)",
	Long:  `Send an email (requires draft id). E.g. "send 12"`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var id int
		id, err := strconv.Atoi(args[0])
		fc.ErrCheck(err, "ID is not an integer")

		ic, err := utils.ImapLogin()
		fc.ErrCheck(err, "error when logging into IMAP server")
		sc, err := utils.SmtpLogin()
		fc.ErrCheck(err, "error when logging into SMTP server")

		defer ic.Logout()
		defer sc.Close()

		draftsFolder := utils.FindMailbox(ic, "\\Drafts", "Drafts")

		subject, body, err := utils.GetMessage(ic, id, draftsFolder)
		fc.ErrCheck(err, "error when getting message details")

		sendPrompt := []*survey.Question{
			{
				Name:     "to",
				Prompt:   &survey.Input{Message: "To (; -separated):"},
				Validate: survey.Required,
			},
		}

		response := struct {
			To   string
			From string
		}{}

		survey.Ask(sendPrompt, &response)

		err = utils.Send(sc, response.To, "vedthiru@hotmail.com", subject, body)
		fc.ErrCheck(err, "error when sending message")
		err = utils.MoveMail(ic, id, draftsFolder, utils.FindMailbox(ic, "\\Sent", "Sent"))
		fc.ErrCheck(err, "error when moving sent draft to Sent folder")

		fc.Success("Sent email.")
	},
}

func init() {
	rootCmd.AddCommand(sendCmd)
}
