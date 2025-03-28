package task

type Controller struct {
	Service IService
}

func NewController(service IService) *Controller {
	return &Controller{service}
}
