package repositories

import (
	"github.com/beego/beego/v2/client/orm"

	"rent-by-owner-api/dtos"
)

type DestinationRepository interface {
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
			OR city % ?
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

func (r *destinationRepository) Autocomplete(
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
			OR city % ?
		ORDER BY
			CASE
				WHEN LOWER(city) = LOWER(?) THEN 1
				WHEN LOWER(city) LIKE LOWER(?) || '%' THEN 2
				WHEN LOWER(country) = LOWER(?) THEN 3
				WHEN LOWER(country) LIKE LOWER(?) || '%' THEN 4
				ELSE 5
			END,
			similarity(city, ?) DESC,
			COALESCE(population, 0) DESC
		LIMIT 5
	`

	var destinations []dtos.DestinationResponse

	_, err := r.orm.Raw(
		sql,
		query, // city partial
		query, // country partial
		query, // city fuzzy
		query, // exact city
		query, // city prefix
		query, // exact country
		query, // country prefix
		query, // similarity ranking
	).QueryRows(&destinations)

	if err != nil {
		return nil, err
	}

	return destinations, nil
}

func (r *destinationRepository) Nearby(
	latitude, longitude, radiusKm float64,
) ([]dtos.NearbyDestinationResponse, error) {
	const sql = `
		SELECT
			id,
			city,
			country,
			COALESCE(population, 0) AS population,
			ST_Y(location::geometry) AS latitude,
			ST_X(location::geometry) AS longitude,
			ST_Distance(
				location,
				ST_SetSRID(
					ST_MakePoint(?, ?),
					4326
				)::geography
			) / 1000 AS distance
		FROM destinations
		WHERE ST_DWithin(
			location,
			ST_SetSRID(
				ST_MakePoint(?, ?),
				4326
			)::geography,
			? * 1000
		)
		ORDER BY distance
	`

	var destinations []dtos.NearbyDestinationResponse

	_, err := r.orm.Raw(
		sql,
		longitude,
		latitude,
		longitude,
		latitude,
		radiusKm,
	).QueryRows(&destinations)

	if err != nil {
		return nil, err
	}

	return destinations, nil
}

func (r *destinationRepository) WithinBounds(
	north float64,
	south float64,
	east float64,
	west float64,
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
		WHERE location && ST_MakeEnvelope(
			?,
			?,
			?,
			?,
			4326
		)::geography
	`

	var destinations []dtos.DestinationResponse

	_, err := r.orm.Raw(
		sql,
		west,
		south,
		east,
		north,
	).QueryRows(&destinations)

	if err != nil {
		return nil, err
	}

	return destinations, nil
}