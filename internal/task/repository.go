package task

import "sync"

type Repository interface {
	Create(request *CreateRequest) (Task, error)
	GetAll() ([]Task, error)
	GetById(id int) (Task, error)
	Update(id int) (Task, error)
	Delete(id int) error
}

type repository struct {
	tasks     map[int]Task
	idCounter int
	locker    sync.RWMutex
}

func NewRepository(tasks map[int]Task) Repository {
	return &repository{tasks: tasks, idCounter: 0}
}
