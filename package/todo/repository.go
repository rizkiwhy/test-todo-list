package todo

import (
	"rizkiwhy/test-todo-list/package/todo/model"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type Repository interface {
	Create(todo model.Todo) (*model.Todo, error)
	GetByID(id string) (*model.Todo, error)
}

type RepositoryImpl struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &RepositoryImpl{DB: db}
}

func (r *RepositoryImpl) Create(todo model.Todo) (*model.Todo, error) {
	result := r.DB.Create(&todo)
	if result.Error != nil {
		log.Error().Err(result.Error).Interface("todo", todo).Msg("[TodoRepository][Create] Failed to create todo")
		return nil, result.Error
	}

	return &todo, nil
}

func (r *RepositoryImpl) GetByID(id string) (*model.Todo, error) {
	var todo model.Todo
	result := r.DB.Where("id = ?", id).First(&todo)
	if result.Error != nil {
		log.Error().Err(result.Error).Str("id", id).Msg("[TodoRepository][GetByID] Failed to get todo by id")
		return nil, result.Error
	}

	return &todo, nil
}
