package task

type Controller struct {
	Service Service
}

func NewController(service Service) *Controller {
	return &Controller{service}
}
