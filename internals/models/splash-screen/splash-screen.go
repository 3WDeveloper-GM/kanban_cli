package splashscreen

import (
	"os/user"

	"github.com/3WDeveloper-GM/kanban_board_new/internals/models/keymap"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	help     help.Model
	Title    textinput.Model
	Confirms bool
	Quitting bool
	Board    tea.Model
	width    int
	height   int
}

func NewModel(title string) *Model {

	screen := &Model{
		help:  help.New(),
		Title: textinput.New(),
	}

	screen.Title.Placeholder = title
	screen.Title.Focus()

	return screen
}

func SplashScreen() *Model {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	strings := "Welcome, " + user.Username
	return NewModel(strings)
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keymap.KeyRing.Quit):
			m.Quitting = true
			return m, tea.Quit
		case key.Matches(msg, keymap.KeyRing.Enter):
			return m.Board.Update(tea.WindowSizeMsg{Width: m.width, Height: m.height})
		}
	}

	if m.Title.Focused() {
		m.Title, cmd = m.Title.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m Model) View() string {
	if m.Quitting {
		return ""
	}
	if !m.Confirms {
		return lipgloss.JoinVertical(
			lipgloss.Center,
			"Please press 'Enter' to start the program",
		)
	} else {
		return m.Board.View()
	}
}
