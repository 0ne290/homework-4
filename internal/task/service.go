package task

type Service interface {
	Create(request CreateRequest) (CreateResponse, error)
	GetAll() ([]GetResponse, error)
	GetById(id int) (GetResponse, error)
	Update(id int) error
	Delete(id int) error
}

type service struct {
	repository repository
}

func NewService(repository Repository) Service {
	return
}
