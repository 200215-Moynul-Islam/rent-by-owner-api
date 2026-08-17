package repositories

import (
	"github.com/beego/beego/v2/client/orm"

	"rent-by-owner-api/dtos"
)

type DestinationRepository interface {
	Search(query string) ([]dtos.DestinationResponse, error)
}

type destinationRepository struct {
	orm orm.Ormer
}

func NewDestinationRepository() DestinationRepository {
	return &destinationRepository{
		orm: orm.NewOrm(),
	}
}

func (r *destinationRepository) Search(
	query string,
) ([]dtos.DestinationResponse, error) {
	const sql = `
		SELECT
			id,
			city,
			country,
			COALESCE(population, 0) AS population,
			ST_Y(location::geometry) AS latitude,
			ST_X(location::geometry) AS longitude
		FROM destinations
		WHERE
			city ILIKE '%' || ? || '%'
			OR country ILIKE '%' || ? || '%'
			OR similarity(city, ?) > 0.5
	`

	var destinations []dtos.DestinationResponse

	_, err := r.orm.Raw(
		sql,
		query,
		query,
		query,
	).QueryRows(&destinations)

	if err != nil {
		return nil, err
	}

	return destinations, nil
}