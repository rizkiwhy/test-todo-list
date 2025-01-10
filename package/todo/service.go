package todo

import "rizkiwhy/test-todo-list/package/todo/model"

type ServiceImpl struct {
	Repository Repository
}

type Service interface {
	Create(todo model.Todo) (*model.Todo, error)
}

func NewService(repository Repository) Service {
	return &ServiceImpl{
		Repository: repository,
	}
}

func (s *ServiceImpl) Create(todo model.Todo) (*model.Todo, error) {
	_, err := s.Repository.Create(todo)
	if err != nil {

		return nil, err
	}

	return &todo, nil
}
