package task

type Service interface {
	Create(request *CreateRequest) (CreateResponse, error)
	GetAll() (GetAllResponse, error)
	GetById(id int) (GetByIdResponse, error)
	Update(id int) (UpdateResponse, error)
	Delete(id int) (DeleteResponse, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{repository}
}
