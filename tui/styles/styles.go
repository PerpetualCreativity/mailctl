package styles

import "github.com/charmbracelet/lipgloss"

// general colors
var (
	notFocused	= lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	focused		= lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	special		= lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}
)

func Focus(s lipgloss.Style, b bool) lipgloss.Style {
	if b { return s.Copy().BorderForeground(focused) } else { return s }
}

// account styling
var (
	// borders taken from https://github.com/charmbracelet/lipgloss/blob/a86f21a0ae430173036c5b6158b0654af447e5a1/example/main.go#L40
	Account = lipgloss.NewStyle().
		Border(lipgloss.Border {
			Top:			"─",
			Bottom:			"─",
			Left:			"│",
			Right:			"│",
			TopLeft:		"╭",
			TopRight:		"╮",
			BottomLeft:		"┴",
			BottomRight:	"┴",
		}, true).
		BorderForeground(notFocused).
		Padding(0, 1)
	ActiveAccount = Account.Copy().
		Border(lipgloss.Border{
			Top:			"─",
			Bottom:			" ",
			Left:			"│",
			Right:			"│",
			TopLeft:		"╭",
			TopRight:		"╮",
			BottomLeft:		"┘",
			BottomRight:	"└",
		}, true).
		BorderForeground(special)
	AccountGap = Account.Copy().
		BorderTop(false).
		BorderLeft(false).
		BorderRight(false).
		BorderForeground(special)
)

// cursor styling
var Cursor = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("6"))

// mailbox styling
var (
	MailboxContainer = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(notFocused)
	Mailbox = lipgloss.NewStyle().
		Margin(0, 1)
)

// message list styling
var (
	MessageListContainer = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(notFocused)
	MessageListSender = lipgloss.NewStyle().
		Padding(0, 1).
		Foreground(lipgloss.Color("4"))
	MessageListSubject = lipgloss.NewStyle()
)

// message styling
var (
	MessageContainer = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(notFocused)
	MessageSender = lipgloss.NewStyle().
		Foreground(lipgloss.Color("4"))
	MessageSubject = lipgloss.NewStyle().
		Foreground(lipgloss.Color("5")).
		Bold(true)
	MessageBody = lipgloss.NewStyle()
)

