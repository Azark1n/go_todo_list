package repository

import (
	"database/sql"
	"go_todo_list/data"
	"go_todo_list/model"
	"time"
)

func Create(t model.Task) (int64, error) {
	res, err := data.Db.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES ($1, $2, $3, $4)", t.Date, t.Title, t.Comment, t.Repeat)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return id, err
}

func FindTasksByDate(date time.Time) ([]model.Task, error) {
	formattedDate := date.Format("20060102")
	rows, err := data.Db.Query("SELECT id, date, title, comment, repeat FROM scheduler WHERE date = ? ORDER BY date", formattedDate)
	return processRows(rows, err)
}

func FindTasksByTitleOrComment(search string) ([]model.Task, error) {
	likePattern := "%" + search + "%"
	rows, err := data.Db.Query("SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date", likePattern, likePattern)
	return processRows(rows, err)
}

func FindAllTasks() ([]model.Task, error) {
	rows, err := data.Db.Query("SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date")
	return processRows(rows, err)
}

func GetTask(id int64) (*model.Task, error) {
	var task model.Task
	err := data.Db.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?", id).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func UpdateTask(task model.Task) error {
	result, err := data.Db.Exec("UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?", task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func DeleteTask(id int64) error {
	_, err := data.Db.Exec("DELETE FROM scheduler WHERE id = ?", id)
	if err != nil {
		return err
	}
	return nil
}

func processRows(rows *sql.Rows, err error) ([]model.Task, error) {
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]model.Task, 0)
	for rows.Next() {
		var task model.Task
		if err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}
