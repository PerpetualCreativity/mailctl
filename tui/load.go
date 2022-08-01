package tui

import (
	"strconv"

	"github.com/PerpetualCreativity/mailctl/utils"
	"github.com/charmbracelet/bubbles/help"
	"github.com/spf13/viper"
)

func load() model {
	m := model{}
	m.height = 40
	m.width = 80

	m.help = help.New()

	for i := 0; i < len(viper.GetStringSlice("accounts")); i++ {
		ad := viper.GetStringMapString("accounts." + strconv.Itoa(i))
		accountName := ad["account_name"]
		fc.ErrNComp(accountName, "", "unable to read `account_name` field in config")
		name := ad["name"]
		fc.ErrNComp(name, "", "unable to read `name` field in config")

		ic, err := utils.ImapLogin(i)
		fc.ErrCheck(err, "error when logging into IMAP server")
		sc, err := utils.SmtpLogin(i)
		fc.ErrCheck(err, "error when logging into SMTP server")

		m.accounts = append(
			m.accounts,
			accountModel{
				accountName: accountName,
				name:        name,
				imapClient:  ic,
				smtpClient:  sc,
				configIndex: index{v: i},
			},
		)
		for j, mailbox := range utils.ListMailboxes(m.accounts[i].imapClient) {
			if mailbox == "INBOX" {
				m.accounts[i].activeMailbox = index{v: j}
			}
			m.accounts[i].mailboxes = append(m.accounts[i].mailboxes, mailboxModel{
				name:           mailbox,
				loadedMessages: -1,
			})
		}
	}
	m.activeAccount = index{v: viper.GetInt("default_account") - 1}
	m = m.loadActiveMailbox(50)

	return m
}

func (m model) loadActiveMailbox(loadNumber int) model {
	activeAccount := m.getActiveAccount()
	activeMailbox := m.getActiveMailbox()
	if loadNumber <= activeMailbox.loadedMessages {
		return m
	}

	activeMailbox.messages = []messageModel{}
	messageList, err := utils.ListMessages(
		activeAccount.imapClient,
		activeMailbox.name,
		uint32(loadNumber),
		uint32(max(1, loadNumber+2-activeMailbox.loadedMessages)),
	)
	m.addErr(err)
	for i := len(messageList) - 1; i >= 0; i-- {
		activeMailbox.messages = append(activeMailbox.messages, messageModel{
			envelope: messageList[i],
		})
	}
	activeMailbox.loadedMessages = loadNumber

	return m
}

func (m model) loadActiveMessageBody() model {
	activeAccount := m.getActiveAccount()
	activeMailbox := m.getActiveMailbox()
	activeMessage := m.getActiveMessage()
	var err error
	_, activeMessage.body, err = utils.GetMessage(
		activeAccount.imapClient,
		int(activeMessage.envelope.SeqNum),
		activeMailbox.name,
	)
	m.addErr(err)
	return m
}
