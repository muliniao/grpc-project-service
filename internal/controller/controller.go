package controller

import (
	"learning/grpc-project-service/internal/service"
)

type (
	Controller interface {
		ProjectController
	}

	controller struct {
		service service.Service
	}
)

func New(service service.Service) Controller {
	return controller{service: service}
}
