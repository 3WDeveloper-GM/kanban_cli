package config

import (
	"database/sql"
	"fmt"
	"os"
	"os/exec"

	"github.com/3WDeveloper-GM/kanban_board_new/internals/data/dbmodels"
	"github.com/3WDeveloper-GM/kanban_board_new/internals/models/column"
	"github.com/3WDeveloper-GM/kanban_board_new/internals/models/column/task"
	"github.com/3WDeveloper-GM/kanban_board_new/internals/models/form"
	"github.com/3WDeveloper-GM/kanban_board_new/internals/models/keymap"
	splashscreen "github.com/3WDeveloper-GM/kanban_board_new/internals/models/splash-screen"
	"github.com/3WDeveloper-GM/kanban_board_new/internals/status"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const margin = 4

type Model struct {
	DSN          string
	Models       dbmodels.TaskModel
	Columns      []column.Model
	form         *form.Model
	help         help.Model
	focused      status.Status
	splashScreen splashscreen.Model
	quitting     bool
	loading      bool
}

func (m Model) Focused() status.Status {
	return m.focused
}

func NewModel() *Model {
	help := help.New()
	help.ShowAll = true
	return &Model{
		help:         help,
		focused:      status.Todo,
		splashScreen: *splashscreen.SplashScreen(),
	}
}

func (m *Model) InitDB(db *sql.DB) {
	m.Models.DB = db
}

func (m *Model) InitBoard(width, height int) {

	m.Columns = []column.Model{
		*column.NewModel(status.Todo, width, height),
		*column.NewModel(status.InProgress, width, height),
		*column.NewModel(status.Done, width, height),
	}

	dummyData, err := m.Models.FetchAll()
	if err != nil {
		panic(err)
	}

	for title, contents := range dummyData {
		m.Columns[title].SetStyle()
		m.Columns[title].Contents.SetItems(contents)
	}

}

type initSplashScreen struct{}

func (m Model) Init() tea.Cmd {
	return func() tea.Msg { return initSplashScreen{} }
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case initSplashScreen:
		m.splashScreen.Board = m
		return m.splashScreen.Update(nil)

	case tea.WindowSizeMsg:

		var cmd tea.Cmd
		var cmds []tea.Cmd
		m.help.Width = msg.Width - margin
		for i := range []status.Status{status.Todo, status.InProgress, status.Done} {
			var res tea.Model
			res, cmd = m.Columns[i].Update(msg)
			m.Columns[i] = res.(column.Model)
			cmds = append(cmds, cmd)
		}
		m.loading = true
		return m, tea.Batch(cmds...)

	case splashscreen.Model:
		if msg.Confirms {
			return m.Update(tea.WindowSizeMsg{})
		}
	case column.RemovedItemMsg:
		m.Models.Delete(msg.ID)

	case column.MoveMsg:
		if m.focused != status.Done {
			err := m.Models.Next(msg.Task)
			if err != nil {
				panic(err)
			}
			return m, m.Columns[m.focused+1].Set(column.APPEND, msg.Task)
		}

	case form.Model:
		InsertedTask := msg.CreateTask(m.focused)
		InsertedTask.ID = msg.ID

		switch {
		case msg.Edition:
			fmt.Println(InsertedTask.ID)
			err := m.Models.Update(InsertedTask)
			if err != nil {
				panic(err)
			}
		case !msg.Edition:
			err := m.Models.Insert(InsertedTask)
			if err != nil {
				panic(err)
			}
		}

		return m, m.Columns[m.focused].Set(msg.Cursor, InsertedTask)

	case tea.KeyMsg:

		if m.Columns[m.focused].Contents.FilterState() == list.Filtering {
			break
		}

		switch {
		case key.Matches(msg, keymap.KeyRing.Quit):
			m.quitting = true
			return m, tea.Quit

		case key.Matches(msg, keymap.KeyRing.Left):
			m.Columns[m.focused].Blur()
			m.focused.Prev()
			m.Columns[m.focused].Focus()
		case key.Matches(msg, keymap.KeyRing.Right):
			m.Columns[m.focused].Blur()
			m.focused.Next()
			m.Columns[m.focused].Focus()

		case key.Matches(msg, keymap.KeyRing.Back):
			m.Columns[m.focused].Contents.SetFilteringEnabled(false)

		case key.Matches(msg, keymap.KeyRing.Search):
			m.Columns[m.focused].Contents.SetFilteringEnabled(true)
			m.Columns[m.focused].Contents.ShowFilter()

		case key.Matches(msg, keymap.KeyRing.New):
			m.form = form.NewDefaultForm()
			m.form.Cursor = column.APPEND
			m.form.Board = m
			return m.form.Update(nil)
		case key.Matches(msg, keymap.KeyRing.Edit):

			selectedTask := m.Columns[m.focused].Contents.SelectedItem().(*task.Task)
			selectedTaskCursor := m.Columns[m.focused].Contents.Cursor()

			m.form = form.NewDefaultForm()
			m.form.ID = selectedTask.ID

			m.form.Edition = true

			m.form.ChangeTitle(selectedTask.Title())
			m.form.Cursor = selectedTaskCursor

			m.form.Board = m
			return m.form.Update(nil)
		}
	}
	res, cmd := m.Columns[m.focused].Update(msg)
	if _, ok := res.(column.Model); ok {
		m.Columns[m.focused] = res.(column.Model)
	} else {
		return res, cmd
	}
	return m, cmd
}

func (m *Model) View() string {
	if m.quitting {
		ClearTerminal()
		return ""
	}

	if !m.loading {
		return "loading"
	}

	board := lipgloss.JoinHorizontal(
		lipgloss.Left,
		m.Columns[status.Todo].View(),
		m.Columns[status.InProgress].View(),
		m.Columns[status.Done].View(),
	)
	return lipgloss.JoinVertical(lipgloss.Left, board, m.help.View(keymap.KeyRing))
}

func ClearTerminal() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}
