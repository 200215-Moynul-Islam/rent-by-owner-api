package utils

import "fmt"

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

	return nil
}