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
		m.addErr(msg.err)
		if msg.modified {
			// TODO: actually edit message
			file, _ := os.Open(msg.filename)
			var contents []byte
			file.Read(contents)
		}
		os.Remove(msg.filename)
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
	switch {
	case key.Matches(msg, m.keys()["Next"]):
		m.focus = focusMailboxes
	case key.Matches(msg, m.keys()["Down"]):
		activeMailbox.activeMessage.incr(len(activeMailbox.messages))
	case key.Matches(msg, m.keys()["Up"]):
		activeMailbox.activeMessage.decr()
	case key.Matches(msg, m.keys()["Open"]):
		isDraftsMailbox := activeMailbox.name == utils.FindMailbox(m.getActiveAccount().imapClient, "\\Drafts", "Drafts")
		m.loadActiveMessageBody()
		return m, openMessage(
			m.getActiveMessage(),
			isDraftsMailbox,
		)
	case key.Matches(msg, m.keys()["New"]):
		isDraftsMailbox := activeMailbox.name == utils.FindMailbox(m.getActiveAccount().imapClient, "\\Drafts", "Drafts")
		if isDraftsMailbox {
			return m, openMessage(
				&messageModel{envelope: utils.Message{Subject: "Subject:\n\n"}},
				true,
			)
		}
	case key.Matches(msg, m.keys()["Move"]):
		mailboxes := m.getActiveAccount().mailboxes
		mailboxNames := make([]string, len(mailboxes))
		for i, mailbox := range mailboxes {
			mailboxNames[i] = mailbox.name
		}
		m.prompts = append(m.prompts, prompt{
			question: "Which mailbox should this message be moved to?",
			choices: mailboxNames,
			activeChoice: m.getActiveAccount().activeMailbox,
			process: func (m model, c string) model {
				err := utils.MoveMail(
					m.getActiveAccount().imapClient,
					int(m.getActiveMessage().envelope.SeqNum),
					m.getActiveMailbox().name,
					c,
				)
				m.addErr(err)
				return m
			},
		})
	case key.Matches(msg, m.keys()["Reply"]):
		isDraftsMailbox := activeMailbox.name == utils.FindMailbox(m.getActiveAccount().imapClient, "\\Drafts", "Drafts")
		if !isDraftsMailbox {
		}
	}
	return m, nil
}
