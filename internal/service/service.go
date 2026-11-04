package service

import (
	"github.com/galiullindo/todo-api/internal/model"
)

type Repository interface {
	GetAll() ([]model.Task, error)
	GetByID(id int64) (model.Task, error)
	Create(task model.Task) (model.Task, error)
}

type service struct {
	repo Repository
}

func New(repo Repository) *service {
	return &service{
		repo: repo,
	}
}

func (s *service) GetAll() ([]model.Task, error) {
	return s.repo.GetAll()
}

func (s *service) GetByID(id int64) (model.Task, error) {
	return s.repo.GetByID(id)
}

func (s *service) Create(title string, status string) (model.Task, error) {
	newTask := model.Task{
		Title:  title,
		Status: status,
	}

	return s.repo.Create(newTask)
}
