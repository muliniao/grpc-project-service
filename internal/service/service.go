package service

import (
	"learning/grpc-project-service/pkg/provider"

	"learning/grpc-project-service/pkg/resource/user"
)

type (
	Service interface {
		provider.Provider
		ProjectService
	}

	service struct {
		provider.AbstractProvider
		//repository   repository.Repository
		userResource user.UserResource
	}
)

func New() Service {
	return &service{
		//repository:   repository,
		userResource: user.New(user.NewConfigFromEnv()),
	}
}

func (s *service) Init() error {
	if err := s.userResource.Init(); err != nil {
		//logging.WithError(err).Errorf("Failed to init user resource")
		return err
	}
	return nil
}
