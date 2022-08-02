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
	focusPrompt
)

// model represents the TUI as a whole.
type model struct {
	accounts      []accountModel
	activeAccount index
	width         int
	height        int
	focus         focusArea
	help          help.Model
	errMessages   []string
	prompts       []prompt
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
// addErr adds err.Error() to errMessages if err != nil
func (m model) addErr(err error) {
	if err != nil {
		m.errMessages = append(m.errMessages, err.Error())
	}
}

// prompt represents a question to the user for
// simple information in the form of a string.
// If choices != nil, the options will be
// displayed and the user will only be able to
// choose among them.
// If choices == nil, the user will provide
// arbitrary input.
// activeChoice holds the current selected
// answer by the user and is a good way to set
// a default answer at the time the prompt is
// first shown.
// Once the user submits their answer, the
// answer will be passed into the specified
// process function. The application model will
// be updated with the result.
type prompt struct {
	question     string
	choices	     []string
	activeChoice index
	process      func (m model, s string) model
}
// accountModel represents an email account,
// and stores the IMAP client and SMTP client
// logged in to the account
type accountModel struct {
	accountName   string
	name          string
	imapClient    *client.Client
	smtpClient    *smtp.Client
	mailboxes     []mailboxModel
	activeMailbox index
	configIndex   index
}
// mailboxModel represents a mailbox and the
// ancillary data required to render it.
type mailboxModel struct {
	name           string
	messages       []messageModel
	loadedMessages int
	activeMessage  index
}
// messageModel is a simple wrapper that links
// a utils.Message envelope with the body.
type messageModel struct {
	envelope utils.Message
	body     string
}
// index is a simple counter type to prevent
// increases/decreases beyond +1
type index struct {
	v int
}

func (i *index) decr() {
	if i.v > 0 {
		i.v = i.v - 1
	}
}
func (i *index) incr(arrayLen int) {
	if i.v+1 < arrayLen {
		i.v = i.v + 1
	}
}

func (m model) Init() tea.Cmd {
	return func() tea.Msg {
		return nil
	}
}
