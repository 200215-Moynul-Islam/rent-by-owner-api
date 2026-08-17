package services

import (
	"rent-by-owner-api/dtos"
	"rent-by-owner-api/repositories"
	"rent-by-owner-api/utils"
)

type DestinationService interface {
	Search(query string) ([]dtos.DestinationResponse, error)
	Autocomplete(query string) ([]dtos.DestinationResponse, error)
	Nearby(latitude, longitude, radiusKm float64) ([]dtos.NearbyDestinationResponse, error)
	WithinBounds(
		north float64,
		south float64,
		east float64,
		west float64,
	) ([]dtos.DestinationResponse, error)
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

func (s *destinationService) Autocomplete(
	query string,
) ([]dtos.DestinationResponse, error) {
	return s.repository.Autocomplete(query)
}

func (s *destinationService) Nearby(
	latitude float64,
	longitude float64,
	radiusKm float64,
) ([]dtos.NearbyDestinationResponse, error) {
	if err := utils.ValidateNearbySearchParameters(
		latitude,
		longitude,
		radiusKm,
	); err != nil {
		return nil, err
	}

	return s.repository.Nearby(
		latitude,
		longitude,
		radiusKm,
	)
}

func (s *destinationService) WithinBounds(
	north float64,
	south float64,
	east float64,
	west float64,
) ([]dtos.DestinationResponse, error) {
	if err := utils.ValidateBounds(
		north,
		south,
		east,
		west,
	); err != nil {
		return nil, err
	}

	return s.repository.WithinBounds(
		north,
		south,
		east,
		west,
	)
}