package task

import (
	"github.com/jackc/pgx/v5"
	"sync"
)

type Repository interface {
	Create(request *CreateRequest) (Task, error)
	GetAll() ([]Task, error)
	GetById(id int) (Task, error)
	Update(id int) (Task, error)
	Delete(id int) error
}

type InMemoryRepository struct {
	tasks     map[int]Task
	idCounter int
	locker    sync.RWMutex
}

type PosgresRepository struct {
	transaction pgx.Tx
}

func newPosgresRepository(transaction pgx.Tx) *PosgresRepository {
	return &PosgresRepository{transaction}
}

func NewInMemoryRepository(tasks map[int]Task) *InMemoryRepository {
	return &InMemoryRepository{tasks: tasks, idCounter: 0}
}
