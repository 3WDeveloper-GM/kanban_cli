package column

import (
	"github.com/3WDeveloper-GM/kanban_board_new/internals/models/column/task"
	"github.com/3WDeveloper-GM/kanban_board_new/internals/models/keymap"
	"github.com/3WDeveloper-GM/kanban_board_new/internals/status"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	focus    bool
	Contents list.Model
	status   status.Status
	width    int
	height   int
	colors   ColorData
}

func (m *Model) Focus() {
	m.focus = true
}

func (m *Model) Blur() {
	m.focus = false
}

func (m *Model) Focused() bool {
	return m.focus
}

func (m Model) Status() status.Status {
	return m.status
}

func NewModel(stats status.Status, width, height int) *Model {
	var focus bool
	if stats == status.Todo {
		focus = true
	}

	colors, err := Palette()
	if err != nil {
		panic(err)
	}

	t_color := lipgloss.AdaptiveColor{Dark: colors.Colors.Color4, Light: colors.Colors.Color12}

	defaultDelegate := list.NewDefaultDelegate()
	defaultDelegate.Styles.SelectedTitle.
		Foreground(t_color).
		BorderLeftForeground(t_color)

	defaultDelegate.Styles.SelectedDesc.
		Foreground(t_color).
		BorderLeftForeground(t_color)

	defaultList := list.New([]list.Item{}, defaultDelegate, width, height)
	defaultList.SetShowHelp(false)
	defaultList.Title = stats.String()
	return &Model{
		focus:    focus,
		Contents: defaultList,
		status:   stats,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.getSize(msg.Width, msg.Height)
		m.Contents.SetSize(msg.Width/margin, msg.Height-margin*2)
	case tea.KeyMsg:

		if m.Contents.FilterState() == list.Filtering {
			break
		}

		switch {

		case key.Matches(msg, keymap.KeyRing.Delete):
			return m, m.deleteCurrent()
		case key.Matches(msg, keymap.KeyRing.Down):
			m.Contents.CursorDown()
			return m, nil
		case key.Matches(msg, keymap.KeyRing.Up):
			m.Contents.CursorUp()
			return m, nil
		case key.Matches(msg, keymap.KeyRing.Enter):
			if m.status != status.Done {
				return m, m.moveToNext()
			}
		case key.Matches(msg, keymap.KeyRing.Back):
			return m, nil
		}
	}

	m.Contents, cmd = m.Contents.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return m.getStyle().Render(m.Contents.View())
}

type RemovedItemMsg struct {
	ID int64
}

func (m *Model) deleteCurrent() tea.Cmd {

	var cmd tea.Cmd
	var ID int64

	if len(m.Contents.VisibleItems()) > 0 {
		task := m.Contents.SelectedItem().(*task.Task)
		ID = task.ID
		m.Contents.RemoveItem(m.Contents.Index())
	}

	m.Contents, cmd = m.Contents.Update(nil)
	return tea.Sequence(cmd, func() tea.Msg { return RemovedItemMsg{ID: ID} })
}

type MoveMsg struct {
	Task *task.Task
}

func (m *Model) moveToNext() tea.Cmd {
	var cur = m.Contents.SelectedItem()
	var dtask *task.Task
	var ok bool

	if dtask, ok = cur.(*task.Task); !ok {
		return nil
	}

	m.Contents.RemoveItem(m.Contents.Index())
	dtask.StatusUp()

	var cmd tea.Cmd
	m.Contents, cmd = m.Contents.Update(nil)

	return tea.Sequence(cmd, func() tea.Msg { return MoveMsg{Task: dtask} })
}

func (m *Model) Set(i int, t *task.Task) tea.Cmd {
	if i != APPEND {
		return m.Contents.SetItem(i, t)
	}
	return m.Contents.InsertItem(APPEND, t)
}
