package tui

import (
	"strings"

	"github.com/PerpetualCreativity/mailctl/tui/styles"
	"github.com/charmbracelet/lipgloss"
)

func max(x int, y int) int { if x < y { return y } else { return x } }
func min0(x int, y int) int {
	if x<0 || y<0 { return 0 }
	if x>y { return y } else { return x }
}
func trim(s string, x int) string {
	if x-1 > len(s) { return s }
	if x-2 < 0 { return "" }
	return s[:x-2]+"…"
}
func cursorPrefix(b bool) string {
	activeCursor := styles.Cursor.Render(">")
	if b { return activeCursor } else { return " " }
}
func wrap(s string, c int) []string {
	if c < 1 {
		return []string{}
	}
	lines := strings.Split(s, "\n")
	wrapped := ""
	for _, line := range lines {
		lineWords := strings.Split(line, " ")
		lineLen := c
		for _, word := range lineWords {
			if l := len(word); l+1 < lineLen {
				lineLen -= l
				wrapped += " " + word
			} else if l >= c {
				lw := len(word)
				for l >= c {
					wrapped += word[lw-l:lw-l+c-1] + "-\n"
					l -= c
				}
				wrapped += word[lw-l:]
			} else {
				wrapped += "\n" + word
			}
		}
		wrapped += "\n"
	}
	return strings.Split(wrapped, "\n")
}

func (m model) View() string {
	doc := strings.Builder{}
	help := m.help.View(m.keys())
	availableHeight := max(0, m.height-lipgloss.Height(help)-1)

	mainView := ""
	if m.focus == focusMessage {
		mainView = renderMessage(
		 	*m.getActiveMessage(),
			m.focus==focusMessage,
		 	m.width,
		 	max(0, availableHeight),
		)
	} else {
		accounts := renderAccounts(m.accounts, m.activeAccount.v, m.width) + "\n"
		doc.WriteString(accounts)

		availableHeight = max(0, availableHeight-lipgloss.Height(accounts))

		activeAccount := *m.getActiveAccount()
		mailboxes := renderMailboxes(
			activeAccount.mailboxes,
			activeAccount.activeMailbox.v,
			m.focus==focusMailboxes,
			availableHeight,
		)

		activeMailbox := *m.getActiveMailbox()
		messageList := renderMessageList(
			activeMailbox.messages,
			activeMailbox.activeMessage.v,
			m.focus==focusMessageList,
			m.width-lipgloss.Width(mailboxes),
			availableHeight,
		)

		mainView = lipgloss.JoinHorizontal(lipgloss.Left, mailboxes, messageList)
	}
	doc.WriteString(mainView)
	// availableHeight = max(0, availableHeight-lipgloss.Height(mainView))
	// doc.WriteString(strings.Repeat("\n", availableHeight))
	doc.WriteString("\n\n")
	doc.WriteString(help)

	return doc.String()
}

func renderAccounts(accounts []accountModel, activeIndex int, width int) string {
	tabs := []string{}
	for i, a := range accounts {
		if i == activeIndex {
			tabs = append(tabs, styles.ActiveAccount.Render(a.name))
		} else {
			tabs = append(tabs, styles.Account.Render(a.name))
		}
	}
	row := lipgloss.JoinHorizontal(lipgloss.Top, tabs...)
	gap := styles.AccountGap.Render(strings.Repeat(" ", max(0, width-lipgloss.Width(row))))
	row = lipgloss.JoinHorizontal(lipgloss.Bottom, row, gap)

	return row
}

func renderMailboxes(
	mailboxes []mailboxModel,
	activeIndex int,
	focus bool,
	availableHeight int,
) string {
	boxes := []string{}
	offset := max(0, activeIndex-availableHeight)
	for i, m := range mailboxes[offset:min0(len(mailboxes),availableHeight+offset)] {
		box := lipgloss.JoinHorizontal(
			lipgloss.Left,
			cursorPrefix(i==activeIndex),
			styles.Mailbox.Render(trim(m.name, 10)),
		)
		boxes = append(boxes, box)
	}
	return styles.Focus(styles.MailboxContainer, focus).Copy().
		Height(availableHeight).
		Render(lipgloss.JoinVertical(lipgloss.Left, boxes...))
}

func renderMessageList(
	messages []messageModel,
	activeIndex int,
	focus bool,
	availableWidth int,
	availableHeight int,
) string {
	messageList := []string{}
	offset := max(0, activeIndex+1-availableHeight)
	for i, m := range messages[offset:min0(len(messages)-1,availableHeight)+offset] {
		sender := styles.MessageListSender.Copy().Width(26).Render(trim(m.envelope.Sender, 25))
		subject := styles.MessageListSubject.Render(trim(m.envelope.Subject, availableWidth-lipgloss.Width(sender)-4))
		messageList = append(messageList, lipgloss.JoinHorizontal(
			lipgloss.Left,
			cursorPrefix(i+offset==activeIndex),
			sender,
			subject,
		))
	}
	return styles.Focus(styles.MessageListContainer, focus).Copy().
		Width(availableWidth).
		Render(lipgloss.JoinVertical(lipgloss.Left, messageList...))
}

func renderMessage(
	message messageModel,
	focus bool,
	availableWidth int,
	availableHeight int,
) string {
	// if len(message.wrappedBody)==0 {
	// 	message.wrappedBody = wrap(message.body, availableWidth)
	// }
	// offset := message.offset.v
	// displayBody := strings.Join(
	// 	message.wrappedBody[offset:min0(len(message.wrappedBody)-1,availableHeight)+offset],
	// 	"\n",
	// )
	messageView := lipgloss.JoinVertical(
		lipgloss.Left,
		styles.MessageSubject.Render(message.envelope.Subject),
		styles.MessageSender.Render(message.envelope.Sender),
		styles.MessageBody.Render(message.body),
	)
	return styles.Focus(styles.MessageContainer, focus).Copy().
		Width(availableWidth).
		Height(availableHeight).
		Render(lipgloss.JoinVertical(lipgloss.Left, messageView))
}
