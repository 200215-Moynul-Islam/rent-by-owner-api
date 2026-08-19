package utils

import (
	"fmt"
	"math"
)

const maxBoundsSizeKm = 50

func ValidateNearbySearchParameters(
	latitude float64,
	longitude float64,
	radiusKm float64,
) error {
	if latitude < -90 || latitude > 90 {
		return fmt.Errorf("latitude must be between -90 and 90")
	}

	if longitude < -180 || longitude > 180 {
		return fmt.Errorf("longitude must be between -180 and 180")
	}

	if radiusKm <= 0 {
		return fmt.Errorf("radius must be greater than 0")
	}

	if radiusKm > 25 {
		return fmt.Errorf("radius must be less than or equal to 25")
	}

	return nil
}

func ValidateBounds(
	north float64,
	south float64,
	east float64,
	west float64,
) error {
	if north < -90 || north > 90 {
		return fmt.Errorf("north must be between -90 and 90")
	}

	if south < -90 || south > 90 {
		return fmt.Errorf("south must be between -90 and 90")
	}

	if east < -180 || east > 180 {
		return fmt.Errorf("east must be between -180 and 180")
	}

	if west < -180 || west > 180 {
		return fmt.Errorf("west must be between -180 and 180")
	}

	if south > north {
		return fmt.Errorf("south cannot be greater than north")
	}

	if west > east {
		return fmt.Errorf("west cannot be greater than east")
	}

	widthKm := calculateDistanceKm(
		north,
		west,
		north,
		east,
	)

	heightKm := calculateDistanceKm(
		south,
		west,
		north,
		west,
	)

	if widthKm > maxBoundsSizeKm {
		return fmt.Errorf(
			"bounds width cannot exceed %d km",
			maxBoundsSizeKm,
		)
	}

	if heightKm > maxBoundsSizeKm {
		return fmt.Errorf(
			"bounds height cannot exceed %d km",
			maxBoundsSizeKm,
		)
	}

	return nil
}

func calculateDistanceKm(
	lat1 float64,
	lon1 float64,
	lat2 float64,
	lon2 float64,
) float64 {
	const earthRadiusKm = 6371.0

	latDifference := (lat2 - lat1) * math.Pi / 180
	lonDifference := (lon2 - lon1) * math.Pi / 180

	lat1Rad := lat1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180

	a :=
		math.Sin(latDifference/2)*math.Sin(latDifference/2) +
			math.Cos(lat1Rad)*
				math.Cos(lat2Rad)*
				math.Sin(lonDifference/2)*
				math.Sin(lonDifference/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadiusKm * c
}