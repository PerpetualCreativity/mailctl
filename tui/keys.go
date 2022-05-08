package tui

import (
	"github.com/charmbracelet/bubbles/key"
)

type keyMap map[string]key.Binding

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{mainKeys["Help"], mainKeys["Quit"]}
}
func (k keyMap) FullHelp() [][]key.Binding {
	list := [][]key.Binding{{},{},{}}
	i := 0
	for _, v := range k {
		list[i%3] = append(list[i%3], v)
		i++
	}
	return list
}

var mainKeys = keyMap{
	"Up": key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "scroll up"),
	),
	"Down": key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	"Next": key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "next section"),
	),
	"Open": key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "open"),
	),
	"New": key.NewBinding(
		key.WithKeys("n"),
		key.WithHelp("n", "new draft"),
	),
	"Move": key.NewBinding(
		key.WithKeys("m"),
		key.WithHelp("m", "move"),
	),
	"Esc": key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("tab", "next section"),
	),
	"Reply": key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "reply"),
	),
	"Help": key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "help"),
	),
	"Quit": key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("^c", "quit"),
	),
}

func (m model) keys() keyMap { return keys[m.focus] }

var coreKeys = keyMap{
	"Help": mainKeys["Help"],
	"Quit": mainKeys["Quit"],
}

var keys = []keyMap{
	// account
	{
		"Up": mainKeys["Up"],
		"Down": mainKeys["Down"],
		"Next": mainKeys["Next"],
	},
	// mailbox
	{
		"Up": mainKeys["Up"],
		"Down": mainKeys["Down"],
		"Next": mainKeys["Next"],
	},
	// messageList
	{
		"Up": mainKeys["Up"],
		"Down": mainKeys["Down"],
		"Next": mainKeys["Next"],
		"Open": mainKeys["Open"],
		"New": mainKeys["New"],
		"Move": mainKeys["Move"],
		"Reply": mainKeys["Reply"],
	},
	// message
	{
		"Up": mainKeys["Up"],
		"Down": mainKeys["Down"],
		"Esc": mainKeys["Esc"],
	},
}

