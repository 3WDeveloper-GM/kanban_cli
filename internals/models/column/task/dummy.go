package task

import (
	"github.com/3WDeveloper-GM/kanban_board_new/internals/status"
	"github.com/charmbracelet/bubbles/list"
)

func DummyData() map[status.Status][]list.Item {
	taskMap := map[status.Status][]list.Item{
		status.Todo: {
			NewItem(status.Todo, "make a model", "with bubbletea"),
			NewItem(status.Todo, "charge my phone", "with a charger"),
			NewItem(status.Todo, "do some programming", "in php"),
			NewItem(status.Todo, "learn javascript", "at dawn"),
		},
		status.InProgress: {
			NewItem(status.InProgress, "do some changes to a blueprint", "using Autocad"),
			NewItem(status.InProgress, "solve some differential equations", "using software"),
			NewItem(status.InProgress, "Make some changes to the bedroom", "install some shelves"),
		},
		status.Done: {
			NewItem(status.Done, "read Berserk", "fuck you, Griffith"),
			NewItem(status.Done, "read Kierkegaard", "the knight of faith..."),
		},
	}
	return taskMap
}
