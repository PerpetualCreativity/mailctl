package tui

import (
	"github.com/PerpetualCreativity/fancyChecks"
	"github.com/PerpetualCreativity/mailctl/utils"
	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-smtp"
)

var fc = fancyChecks.New("", "", "", "Error")

func Start() {
	p := tea.NewProgram(load(), tea.WithAltScreen())
	fc.ErrCheck(p.Start(), "Could not start interface")
}

type focusArea int

const (
	focusAccounts = iota
	focusMailboxes
	focusMessageList
	focusMessage
)

type model struct {
	accounts		[]accountModel
	activeAccount	index
	width			int
	height			int
	focus			focusArea
	help			help.Model
}
func (m model) getActiveAccount() *accountModel {
	return &m.accounts[m.activeAccount.v]
}
func (m model) getActiveMailbox() *mailboxModel {
	activeAccount := m.getActiveAccount()
	return &activeAccount.mailboxes[activeAccount.activeMailbox.v]
}
func (m model) getActiveMessage() *messageModel {
	activeMailbox := m.getActiveMailbox()
	return &activeMailbox.messages[activeMailbox.activeMessage.v]
}
type accountModel struct {
	accountName		string
	name			string
	imapClient		*client.Client
	smtpClient		*smtp.Client
	mailboxes		[]mailboxModel
	activeMailbox	index
	configIndex		index
}
type mailboxModel struct {
	name			string
	messages		[]messageModel
	loadedMessages	int
	activeMessage	index
}
type messageModel struct {
	envelope		utils.Message
	body			string
}
type index struct {
	v int
}
func (i *index) decr() { if i.v > 0 { i.v=i.v-1 } }
func (i *index) incr(arrayLen int) { if i.v+1 < arrayLen { i.v=i.v+1 } }

func (m model) Init() tea.Cmd {
	return func() tea.Msg {
		return nil
	}
}

