package cmd

import (
	"github.com/PerpetualCreativity/mailctl/utils"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// viewCmd represents the view command
var viewCmd = &cobra.Command{
	Use:   "view id [folder]",
	Short: "View an email or draft. Requires an id.",
	Long: `Use this command to view the email or draft
associated with the id you gave. If you want to view
an email in a folder other than the inbox, add its name
after the id.`,
	Args: cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		var id int
		id, err := strconv.Atoi(args[0])
		fc.ErrCheck(err, "ID is not an integer")

		c, err := utils.ImapLogin()
		fc.ErrCheck(err, "error when logging into IMAP server")
		defer c.Logout()

		folder := "INBOX"
		if len(args) > 1 {
			folder = args[1]
		}

		subject, body, err := utils.GetMessage(c, id, folder)
		fc.ErrCheck(err, "error when getting message details")

		f, err := os.CreateTemp("", "*.md")
		fc.ErrCheck(err, "Could not create temporary file")
		defer os.Remove(f.Name())
		_, err = f.WriteString(body)
		fc.ErrCheck(err, "Could not display body of message")

		fc.Neutral("Displaying " + subject)
		pagerCmd := strings.Split(os.ExpandEnv("$PAGER"), " ")
		pager := exec.Command(pagerCmd[0], f.Name())
		pager.Stdin = os.Stdin
		pager.Stdout = os.Stdout
		err = pager.Run()
		fc.ErrCheck(err, "Pager quit")
		time.Sleep(3000)
	},
}

func init() {
	rootCmd.AddCommand(viewCmd)
}
