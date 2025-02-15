package api

import (
	"github.com/go-playground/validator/v10"
	"merch/internal/service"
)

type Api struct {
	service *service.Service
}

func New(service *service.Service) *Api {
	return &Api{
		service: service,
	}
}

func validate(request interface{}) error {
	valid := validator.New(validator.WithRequiredStructEnabled())

	err := valid.Struct(request)
	if err != nil {
		return err
	}
	return nil
}
