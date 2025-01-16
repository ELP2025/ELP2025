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

func compute_new_state(S_n float32, S_m float32, B float32, S float32, K float32) uint8 {
	birth := sigmoid(S_m, B, K)
	survival := sigmoid(S_m, B, K)

	calcul_float := S_n*birth + (1-S_n)*survival
	val_couleur := calcul_float * 255
	roundedUp := math.Ceil(val_couleur)
	intNumber := uint8(roundedUp)
	return intNumber
}

func updateLine(pixels [][][]uint8, y int) [][]uint8 {
	newPixels := make([][]uint8, width)
	for x := range newPixels {
		S_n, S_m, B, S, K := func_gabi()
		new_color := compute_new_state(S_n, S_m, B, S, K)
		newPixels[x] = []uint8{new_color, new_color, new_color}
	}
	return newPixels
}

func UpdateGrid(pixels [][][]uint8) [][][]uint8 {
	newPixels := make([][][]uint8, height)
	done := make(chan int, height) // Channel to synchronize goroutines

	for i := 0; i < height; i++ {
		i := i // Capture the loop variable
		go func() {
			newPixels[i] = updateLine(pixels, i)
			done <- i // Signal completion
		}()
	}

	// Wait for all goroutines to finish
	for i := 0; i < height; i++ {
		<-done
	}
	return newPixels
}
