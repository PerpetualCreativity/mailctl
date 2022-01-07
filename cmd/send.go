package cmd

import (
	"mailctl/utils"
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
		utils.ErrCheck(err, "ID is not an integer")

		ic := utils.ImapLogin()
		sc := utils.SmtpLogin()

		defer ic.Logout()
		defer sc.Close()

		draftsFolder := utils.FindMailbox(ic, "\\Drafts", "Drafts")

		subject, body := utils.GetMessage(ic, id, draftsFolder)

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

		utils.Send(sc, response.To, "vedthiru@hotmail.com", subject, body)

		utils.MoveMail(ic, id, draftsFolder, utils.FindMailbox(ic, "\\Sent", "Sent"))

		utils.Success("Sent email.")
	},
}

func init() {
	rootCmd.AddCommand(sendCmd)
}
