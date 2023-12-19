package keymap

import "github.com/charmbracelet/bubbles/key"

type Keymap struct {
	New    key.Binding
	Edit   key.Binding
	Delete key.Binding
	Up     key.Binding
	Down   key.Binding
	Right  key.Binding
	Left   key.Binding
	Help   key.Binding
	Quit   key.Binding
	Back   key.Binding
	Enter  key.Binding
	Search key.Binding
}

var KeyRing = Keymap{
	New: key.NewBinding(
		key.WithKeys("n"),
		key.WithHelp("n", "new"),
	),
	Edit: key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "edit"),
	),
	Delete: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "delete"),
	),
	Up: key.NewBinding(
		key.WithKeys("k"),
		key.WithHelp("↑/k", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("j"),
		key.WithHelp("↓/j", "down"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "move right"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/l", "move left"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q/ctrl+c", "quit"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "enter"),
	),
	Search: key.NewBinding(
		key.WithKeys("S", "/"),
		key.WithHelp("S", "save"),
	),
}

func (k Keymap) ShortHelp() []key.Binding {
	return []key.Binding{k.Quit, k.Help}
}

func (k Keymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right},
		{k.Help, k.Quit, k.Delete, k.Back},
	}
}
