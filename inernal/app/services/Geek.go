package services

import (
	"GeekReward/inernal/app/models/tables"
	"GeekReward/inernal/app/repositories"
)

type GeekService interface {
	GetTopGeeks(limit int) ([]tables.User, error)
	GetGeekByID(id uint) (*tables.User, error)
}

type geekService struct {
	geekRepo repositories.GeekRepository
}

func NewGeekService(geekRepo repositories.GeekRepository) GeekService {
	return &geekService{geekRepo: geekRepo}
}

func (s *geekService) GetTopGeeks(limit int) ([]tables.User, error) {
	return s.geekRepo.GetTopGeeks(limit)
}

func (s *geekService) GetGeekByID(id uint) (*tables.User, error) {
	return s.geekRepo.GetGeekByID(id)
}
