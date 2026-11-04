package repository

import (
	"database/sql"
	"log"

	"github.com/galiullindo/todo-api/internal/model"
)

type repository struct {
	db *sql.DB
}

func New(db *sql.DB) *repository {
	return &repository{
		db: db,
	}
}

func (r *repository) GetAll() ([]model.Task, error) {
	rows, err := r.db.Query("SELECT id, title, status FROM tasks")
	if err != nil {
		log.Println("repository:", err)
		return nil, err
	}
	defer rows.Close()
	tasks := make([]model.Task, 0)

	for rows.Next() {
		var task model.Task
		if err := rows.Scan(&task.ID, &task.Title, &task.Status); err != nil {
			log.Println("repository:", err)
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if err := rows.Close(); err != nil {
		log.Println("repository:", err)
		return nil, err
	}
	if err := rows.Err(); err != nil {
		log.Println("repository:", err)
		return nil, err
	}

	return tasks, nil
}

func (r *repository) GetByID(id int64) (model.Task, error) {
	var task model.Task
	row := r.db.QueryRow("SELECT id, title, status FROM tasks WHERE id=?", id)

	err := row.Scan(&task.ID, &task.Title, &task.Status)
	if err != nil {
		log.Println("repository:", err)
		return task, err
	}

	return task, nil
}

func (r *repository) Create(task model.Task) (model.Task, error) {
	result, err := r.db.Exec("INSERT INTO tasks (title, status) VALUES (?, ?)", task.Title, task.Status)
	if err != nil {
		log.Println("repository:", err)
		return task, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Println("repository:", err)
		return task, err
	}
	task.ID = id

	return task, nil
}
