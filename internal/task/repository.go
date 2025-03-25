package task

type Repository interface {
	Create(request CreateRequest) (Task, error)
	GetAll() ([]Task, error)
	GetById(id int) (Task, error)
	Update(id int) error
	Delete(id int) error
}

type repository struct {
	tasks map[int]Task
}

func NewRepository(tasks map[int]Task) Repository {
	return &repository{tasks}
}

func (r *repository) Create(request CreateRequest) (Task, error) {

}

func (r *repository) GetAll() ([]Task, error) {

}

func (r *repository) GetById(id int) (Task, error) {

}

func (r *repository) Update(id int) error {

}

func (r *repository) Delete(id int) error {

}
