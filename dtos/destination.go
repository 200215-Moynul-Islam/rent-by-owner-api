package dtos

type DestinationResponse struct {
	Id         int64   `json:"id"`
	City       string  `json:"city"`
	Country    string  `json:"country"`
	Population int64   `json:"population"`
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
}

type NearbyDestinationResponse struct {
	DestinationResponse
	Distance   float64 `json:"distance"`
}