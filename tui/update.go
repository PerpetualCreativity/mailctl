package tui

import (
	"os"
	"os/exec"

	"github.com/PerpetualCreativity/mailctl/utils"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type editorEnd struct {
	err      error
	filename string
	modified bool
	message  *messageModel
}

func openMessage(m *messageModel, modifiable bool) tea.Cmd {
	temp, _ := os.CreateTemp("", "*.txt")
	defer os.Remove(temp.Name())
	temp.WriteString(m.body)
	if modifiable {
		temp.Chmod(0660)
	} else {
		temp.Chmod(0440)
	}
	temp.Close()
	c := exec.Command(os.Getenv("EDITOR"), temp.Name())
	return tea.ExecProcess(c, func(err error) tea.Msg {
		return editorEnd{
			err:      err,
			filename: temp.Name(),
			modified: modifiable,
			message:  m,
		}
	})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, coreKeys["Quit"]):
			return m, tea.Quit
		case key.Matches(msg, coreKeys["Help"]):
			m.help.ShowAll = !m.help.ShowAll
			return m, nil
		}
		switch m.focus {
		case focusAccounts:
			return m.accountsUpdate(msg)
		case focusMailboxes:
			return m.mailboxesUpdate(msg)
		case focusMessageList:
			return m.messageListUpdate(msg)
		}
	case tea.WindowSizeMsg:
		m.height = msg.Height - 2
		m.width = msg.Width - 2
		m.help.Width = msg.Width
	case editorEnd:
		fc.ErrCheck(msg.err, "Error from editor")
		// TODO: actually edit message
		if msg.modified {
			file, _ := os.Open(msg.filename)
			contents := []byte{}
			file.Read(contents)
		}
		return m, tea.EnterAltScreen
	}
	return m, nil
}

func (m model) accountsUpdate(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys()["Next"]):
		m.focus = focusMailboxes
	case key.Matches(msg, m.keys()["Down"]):
		m.activeAccount.incr(len(m.accounts))
	case key.Matches(msg, m.keys()["Up"]):
		m.activeAccount.decr()
	}
	return m, nil
}

func (m model) mailboxesUpdate(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	activeAccount := m.getActiveAccount()
	switch {
	case key.Matches(msg, m.keys()["Next"]):
		m.focus = focusMessageList
	case key.Matches(msg, m.keys()["Down"]):
		activeAccount.activeMailbox.incr(len(activeAccount.mailboxes))
		m.loadActiveMailbox(50)
	case key.Matches(msg, m.keys()["Up"]):
		activeAccount.activeMailbox.decr()
		m.loadActiveMailbox(50)
	}
	return m, nil
}

func (m model) messageListUpdate(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	activeMailbox := m.getActiveMailbox()
	isDraftsMailbox := m.getActiveMailbox().name == utils.FindMailbox(m.getActiveAccount().imapClient, "\\Drafts", "Drafts")
	switch {
	case key.Matches(msg, m.keys()["Next"]):
		m.focus = focusMailboxes
	case key.Matches(msg, m.keys()["Down"]):
		activeMailbox.activeMessage.incr(len(activeMailbox.messages))
	case key.Matches(msg, m.keys()["Up"]):
		activeMailbox.activeMessage.decr()
	case key.Matches(msg, m.keys()["Open"]):
		m.loadActiveMessageBody()
		return m, openMessage(
			m.getActiveMessage(),
			isDraftsMailbox,
		)
	case key.Matches(msg, m.keys()["New"]):
		if isDraftsMailbox {
			return m, openMessage(
				&messageModel{envelope: utils.Message{Subject: "Subject:\n\n"}},
				true,
			)
		}
	case key.Matches(msg, m.keys()["Move"]):
	case key.Matches(msg, m.keys()["Reply"]):
		if !isDraftsMailbox {
		}
	}
	return m, nil
}
