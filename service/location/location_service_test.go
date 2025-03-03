package location

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFindDistance_SameLocation(t *testing.T) {
	result := FindDistanceInMeters(0, 0, 0, 0)
	assert.Equal(t, 0.0, result, "Distance between the same points should be 0")
}

func TestFindDistance_OneDegreeLatitude(t *testing.T) {
	result := FindDistanceInMeters(0, 0, 1, 0)
	expected := 111000.0
	assert.InDelta(t, expected, result, 1000, "Distance for 1° latitude should be ~111 km")
}

func TestFindDistanceOf3935KM(t *testing.T) {
	result := FindDistanceInMeters(40.7128, -74.0060, 34.0522, -118.2437)
	expected := 3935000.0
	assert.InDelta(t, expected, result, 10000, "Distance should be ~3935 km")
}

func TestFindDistanceOf222KM(t *testing.T) {
	result := FindDistanceInMeters(-1, 0, 1, 0)
	expected := 222000.0 // Roughly 222 km
	assert.InDelta(t, expected, result, 1000, "Distance across the equator should be ~222 km")
}
