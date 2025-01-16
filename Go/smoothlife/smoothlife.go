package smoothlife

import (
	"math"
	"math/rand"
)

var (
	width  = 1000
	height = 1000
)

// genere des pixels avec une couleur random
func GenerateRandomPixels(grid_width int, grid_height int) [][][]uint8 {
	width = grid_width
	height = grid_height
	nestedPixels := make([][][]uint8, height)
	for y := range nestedPixels {
		nestedPixels[y] = make([][]uint8, width)
		for x := range nestedPixels[y] {
			aliveness := rand.Float32() / 2
			nestedPixels[y][x] = []uint8{
				uint8(255 * aliveness), // R
				uint8(255 * aliveness), // G
				uint8(255 * aliveness), // B
			}
		}
	}
	return nestedPixels
}

func sigmoid(x float32, threshold float32, steepness float32) float32 {
	return 1 / (1 + math.Exp(-steepness*(x-threshold)))
}
