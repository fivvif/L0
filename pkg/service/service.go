package service

import (
	"L0"
	"L0/pkg/repository"
)

type Cache interface {
	NewCache() error
	AddOrder(orderUID string, order L0.Order)
	GetCache(uid string) (L0.Order, error)
}

type Service struct {
	Cache
}

func NewService(repo *repository.Repository) (*Service, error) {
	cache, err := NewCacheService(repo)
	if err != nil {
		return nil, err
	}
	return &Service{Cache: cache}, nil
}
