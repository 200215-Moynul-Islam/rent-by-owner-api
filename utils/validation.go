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

	return nil
}