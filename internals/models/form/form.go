package form

import (
	"github.com/3WDeveloper-GM/kanban_board_new/internals/models/column/task"
	"github.com/3WDeveloper-GM/kanban_board_new/internals/models/keymap"
	"github.com/3WDeveloper-GM/kanban_board_new/internals/status"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	help        help.Model
	title       textinput.Model
	description textarea.Model
	Edition     bool
	Cursor      int
	quitting    bool
	ID          int64
	Board       tea.Model
}

func (m *Model) ChangeTitle(title string) {
	m.title.Placeholder = title
}

func NewModel(title, description string) *Model {
	form := &Model{
		help:        help.New(),
		title:       textinput.New(),
		description: textarea.New(),
	}
	form.title.Placeholder = title
	form.description.Placeholder = description
	form.title.Focus()
	return form
}

func NewDefaultForm() *Model {
	return NewModel("task name", "")
}

func (m Model) CreateTask(status status.Status) *task.Task {
	return task.NewItem(
		status,
		m.title.Value(),
		m.description.Value(),
	)
}

func (f Model) Init() tea.Cmd {
	return nil
}

func (f Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keymap.KeyRing.Quit):
			f.quitting = true
			return f, tea.Quit
		case key.Matches(msg, keymap.KeyRing.Enter):
			if f.title.Focused() {
				f.title.Blur()
				f.description.Focus()
				return f, textarea.Blink
			}
			// Return the completed form as a message.
			return f.Board.Update(f)
		case key.Matches(msg, keymap.KeyRing.Back):
			return f.Board.Update(nil)
		}
	}
	if f.title.Focused() {
		f.title, cmd = f.title.Update(msg)
		return f, cmd
	}
	f.description, cmd = f.description.Update(msg)
	return f, cmd
}

func (f Model) View() string {
	if f.quitting {
		return ""
	}
	if f.Edition {
		return lipgloss.JoinVertical(
			lipgloss.Left,
			"Edit existing task",
			f.title.View(),
			f.description.View(),
			f.help.View(keymap.KeyRing))
	} else {
		return lipgloss.JoinVertical(
			lipgloss.Left,
			"Create a new task",
			f.title.View(),
			f.description.View(),
			f.help.View(keymap.KeyRing))
	}
}
