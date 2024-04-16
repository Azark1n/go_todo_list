package service

import (
	"database/sql"
	stderrors "errors"
	"go_todo_list/errors"
	"go_todo_list/model"
	"go_todo_list/repository"
	"time"
)

func CreateTask(task model.Task) (int64, error) {
	if task.Title == "" {
		return 0, &errors.ValidationError{Message: "title is not correct"}
	}

	now := time.Now()
	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	if task.Date == "" {
		task.Date = now.Format("20060102")
	} else {
		date, err := time.Parse("20060102", task.Date)
		if err != nil {
			return 0, &errors.ValidationError{Message: "date is not correct"}
		}
		if date.Before(now) {
			if task.Repeat == "" {
				task.Date = now.Format("20060102")
			} else {
				task.Date, err = NextDate(now, date, task.Repeat)
				if err != nil {
					return 0, &errors.ValidationError{Message: "repeat is not correct"}
				}
			}
		}
	}
	return repository.Create(task)
}

func GetTask(id int64) (*model.Task, error) {
	task, err := repository.GetTask(id)
	if stderrors.Is(err, sql.ErrNoRows) {
		return nil, &errors.NotFoundError{Message: "task not found"}
	}
	return task, err
}

func UpdateTask(task model.Task) error {
	if task.ID == "" || task.Title == "" || task.Date == "" || task.Comment == "" {
		return &errors.ValidationError{Message: "required data is missing"}
	}
	err := repository.UpdateTask(task)
	if stderrors.Is(err, sql.ErrNoRows) {
		return &errors.NotFoundError{Message: "task not found"}
	}
	return err
}

func DeleteTask(id int64) error {
	err := repository.DeleteTask(id)
	if stderrors.Is(err, sql.ErrNoRows) {
		return &errors.NotFoundError{Message: "task not found"}
	}
	return err
}

func MarkDoneTask(id int64) error {
	task, err := repository.GetTask(id)
	if stderrors.Is(err, sql.ErrNoRows) {
		return &errors.NotFoundError{Message: "task not found"}
	}

	if task.Repeat == "" {
		return repository.DeleteTask(id)
	} else {
		now, _ := time.Parse("20060102", task.Date)
		nextDate, err := NextDate(now, now, task.Repeat)
		if err != nil {
			return err
		}
		task.Date = nextDate
		return repository.UpdateTask(*task)
	}
}

func FindTasks(search string) ([]model.Task, error) {
	if search == "" {
		return repository.FindAllTasks()
	}
	if date, err := time.Parse("02.01.2006", search); err == nil {
		return repository.FindTasksByDate(date)
	}
	return repository.FindTasksByTitleOrComment(search)
}
