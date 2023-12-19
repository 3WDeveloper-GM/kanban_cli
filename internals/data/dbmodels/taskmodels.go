package dbmodels

import (
	"database/sql"
	"time"

	"github.com/3WDeveloper-GM/kanban_board_new/internals/models/column/task"
	"github.com/3WDeveloper-GM/kanban_board_new/internals/status"
	"github.com/charmbracelet/bubbles/list"
)

type TaskModel struct {
	DB *sql.DB
}

func NewModel(db *sql.DB) TaskModel {
	return TaskModel{
		DB: db,
	}
}

func (t *TaskModel) Insert(task *task.Task) error {
	var query = `
	 INSERT INTO tasks (title, description, status)
	 VALUES ($1, $2, $3)
	 RETURNING id, created_at, version
	 `

	args := []interface{}{task.Title(), task.Description(), task.GetStatus().String()}

	return t.DB.QueryRow(query, args...).Scan(&task.ID, &task.CreatedAt, &task.Version)
}

func (t *TaskModel) FetchAll() (map[status.Status][]list.Item, error) {
	var query = `
	SELECT id, created_at, title, description
	FROM tasks
	WHERE status = $1
	ORDER BY id
	`

	initialMap := make(map[status.Status][]list.Item)

	for _, boardTitle := range []status.Status{status.Todo, status.InProgress, status.Done} {
		result, err := t.DB.Query(query, boardTitle.String())

		if err != nil {
			return nil, err
		}

		for result.Next() {
			var ID int64
			var time time.Time
			var title string
			var description string

			err := result.Scan(&ID, &time, &title, &description)
			if err != nil {
				return nil, err
			}

			newtask := task.NewItem(boardTitle, title, description)
			newtask.ID = ID
			newtask.CreatedAt = time

			initialMap[boardTitle] = append(initialMap[boardTitle], newtask)
		}

	}
	return initialMap, nil
}

func (t *TaskModel) Delete(id int64) error {
	var query = `
	DELETE FROM tasks
	WHERE id = $1
	`

	_, err := t.DB.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}

func (t *TaskModel) Update(someTask *task.Task) error {
	var query = `
	UPDATE tasks
	SET title = $1, description = $2, version = version +1
	WHERE id = $3;
	`

	args := []interface{}{someTask.Title(), someTask.Description(), someTask.ID}
	_, err := t.DB.Exec(query, args...)

	return err
}

func (t *TaskModel) Next(someTask *task.Task) error {
	var query = `
	UPDATE tasks
	SET status = $1, version = version +1 
	WHERE id = $2
	`

	_, err := t.DB.Exec(query, someTask.GetStatus().String(), someTask.ID)
	return err
}
