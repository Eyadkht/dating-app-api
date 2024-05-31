package utils

import (
	"math"
	"math/rand"
)

// Generate random latitude and longitude
func GenerateRandomLatLong() (float64, float64) {
	lat := -90.0 + rand.Float64()*180.0
	lon := -180.0 + rand.Float64()*360.0
	return lat, lon
}

func CalculateDistance(lat1, lon1, lat2, lon2 float64) float64 {

	if lat1 == 0 || lon1 == 0 || lat2 == 0 || lon2 == 0 {
		// Returning large distance in case the coordinates are wrong
		// prioritising users with known locations
		return 10000000000
	}

	// Convert degrees to radians
	lat1 = lat1 * math.Pi / 180
	lon1 = lon1 * math.Pi / 180
	lat2 = lat2 * math.Pi / 180
	lon2 = lon2 * math.Pi / 180

	// Haversine formula
	// https://www.geeksforgeeks.org/haversine-formula-to-find-distance-between-two-points-on-a-sphere/
	dlon := lon2 - lon1
	dlat := lat2 - lat1
	a := math.Pow(math.Sin(dlat/2), 2) + math.Cos(lat1)*math.Cos(lat2)*math.Pow(math.Sin(dlon/2), 2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	// Earth's radius in kilometers
	distance := 6371 * c

	return distance
}
