package utils

import (
	"sort"
	"strings"

	"rent-by-owner-api/dtos"
)

func ReRankDestinations(
	query string,
	destinations []dtos.DestinationResponse,
) []dtos.DestinationResponse {
	normalizedQuery := strings.ToLower(
		strings.TrimSpace(query),
	)

	sort.SliceStable(destinations, func(i, j int) bool {
		scoreI := calculateDestinationScore(
			normalizedQuery,
			destinations[i],
		)

		scoreJ := calculateDestinationScore(
			normalizedQuery,
			destinations[j],
		)

		return scoreI > scoreJ
	})

	return destinations
}

func calculateDestinationScore(
	query string,
	destination dtos.DestinationResponse,
) float64 {
	city := strings.ToLower(
		strings.TrimSpace(destination.City),
	)

	country := strings.ToLower(
		strings.TrimSpace(destination.Country),
	)

	var score float64

	// Exact city match.
	if city == query {
		score += 100
	}

	// City prefix match.
	if strings.HasPrefix(city, query) {
		score += 70
	}

	// City partial match.
	if strings.Contains(city, query) {
		score += 30
	}

	// Exact country match.
	if country == query {
		score += 50
	}

	// Country prefix match.
	if strings.HasPrefix(country, query) {
		score += 30
	}

	// Population is only a secondary signal.
	score += populationScore(destination.Population)

	return score
}

func populationScore(population int64) float64 {
	if population <= 0 {
		return 0
	}

	switch {
	case population >= 10_000_000:
		return 10
	case population >= 1_000_000:
		return 8
	case population >= 100_000:
		return 5
	case population >= 10_000:
		return 2
	default:
		return 1
	}
}