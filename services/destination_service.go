package services

import (
	"rent-by-owner-api/dtos"
	"rent-by-owner-api/repositories"
)

type DestinationService interface {
	Search(query string) ([]dtos.DestinationResponse, error)
}

type destinationService struct {
	repository repositories.DestinationRepository
}

func NewDestinationService(
	repository repositories.DestinationRepository,
) DestinationService {
	return &destinationService{
		repository: repository,
	}
}

func (s *destinationService) Search(
	query string,
) ([]dtos.DestinationResponse, error) {
	return s.repository.Search(query)
}