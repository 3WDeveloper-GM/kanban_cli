package task

import (
	"time"

	"github.com/3WDeveloper-GM/kanban_board_new/internals/status"
)

type Task struct {
	status      status.Status
	title       string
	description string
	CreatedAt   time.Time
	Version     int32
	ID          int64
}

func NewItem(status status.Status, title, desc string) *Task {
	newItem := &Task{
		status:      status,
		title:       title,
		description: desc,
	}

	return newItem
}

func (t *Task) GetTitle() *string {
	return &t.title
}

func (t *Task) GetDescription() *string {
	return &t.description
}

func (t Task) FilterValue() string {
	return t.title
}

func (t Task) Title() string {
	return t.title
}

func (t Task) Description() string {
	description := t.description
	return description
}

func (t *Task) GetStatus() *status.Status {
	return &t.status
}

func (t *Task) StatusUp() {
	if t.status != status.Done {
		t.status++
	}
}
