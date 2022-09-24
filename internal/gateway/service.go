package gateway

import (
	"context"
)

type (
	Router interface {
		Forward(ctx context.Context, body interface{}) (interface{}, error)
	}

	Service struct {
		mapper map[string]Router
	}
)

func NewService(mapper map[string]Router) *Service {
	return &Service{
		mapper: mapper,
	}
}

func (s *Service) ResolveService(provider string) Router {
	return s.mapper[provider]
}

func (s *Service) Forward(ctx context.Context, body interface{}, serviceProvider Router) (interface{}, error) {
	res, err := serviceProvider.Forward(ctx, body)
	if err != nil {
		return nil, err
	}
	return res, nil
}
