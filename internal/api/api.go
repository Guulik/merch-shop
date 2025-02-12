package api

import (
	"merch/internal/service"
)

type Api struct {
	//log

	service *service.Service
}

func New(service *service.Service) *Api {
	return &Api{
		service: service,
	}
}
