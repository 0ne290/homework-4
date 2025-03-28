package task

type IService interface {
	Create(request *CreateRequest) (CreateResponse, error)
	GetAll() (GetAllResponse, error)
	GetById(id int) (GetByIdResponse, error)
	Update(id int) (UpdateResponse, error)
	Delete(id int) (DeleteResponse, error)
}

type Service struct {
	unitOfWork UnitOfWork
}

func NewService(unitOfWork UnitOfWork) *Service {
	return &Service{unitOfWork}
}
