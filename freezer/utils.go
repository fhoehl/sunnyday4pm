package freezer

import (
	"math/rand"
)

// Remap value from [min, max] to [toMin, toMax]
func mapValue(value, min, max, toMin, toMax float32) float32 {
	return (value-min)*(toMax-toMin)/(max-min) + toMin
}

// Shuffle an array of string
func shuffle(a []string) {
	for i := range a {
		j := r.Intn(i + 1)
		a[i], a[j] = a[j], a[i]
	}
}

// Return a random number between min and max
func randomBetween(min, max float32) float32 {
	if min > max {
		min, max = max, min
	}

	return rand.Float32()*(max-min) + min
}
