package location

import (
	"math"
)

func FindDistanceInMeters(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadius = 6371000

	lat1, lon1 = toRadians(lat1), toRadians(lon1)
	lat2, lon2 = toRadians(lat2), toRadians(lon2)

	deltaLat := lat2 - lat1
	deltaLon := lon2 - lon1

	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1)*math.Cos(lat2)*math.Sin(deltaLon/2)*math.Sin(deltaLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadius * c
}

func toRadians(deg float64) float64 {
	return deg * (math.Pi / 180)
}
