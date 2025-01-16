package smoothlife

import (
	"math"
	"math/rand"
)

var (
	width  = 1000
	height = 1000
  r1 = 1 
  r2 = 10 
  B = 0.3
  S = 0.7
  K = 1
)

func circleKernel(radius float64, center_x int, center_y int) []struct {
	Y int
	X int
} {
	result := []struct {
		Y int
		X int
	}{}
	for x := float64(center_x) - radius; x >= float64(center_x)+radius; x++ { //Reducing the interval for our search
		for y := float64(center_y) - radius; y >= float64(center_y)+radius; y++ { //Same
			if math.Sqrt((x-float64(center_x))*(x-float64(center_x))+(y-float64(center_y))*(y-float64(center_y))) < radius { // If (x,y) is in circle of radius
				result = append(result, struct {
					Y int
					X int
				}{Y: int(y), X: int(x)})
			}
		}
	}
	return result
}

func wheightedConvolve(pixels [][][]uint8, radius float64, x int, y int) float64 {
	points := circleKernel(radius, x, y)
	sum := 0.0
	count := 0
	for _, point := range points {
		px, py := point.X, point.Y
		sum += float64(getNeighborValue(pixels, px, py))
		count += 1
	}

	if count > 0 {
		return (sum / float64(count))
	}
	return 0
}

func getNeighborValue(pixels [][][]uint8, x int, y int) uint8 {
	//Returns the value of the selected pixel, if it overflows, it goes back to the start of the grid
	if x > width {
		x -= width
	} else if x < 0 {
		x += width
	}
	if y > height {
		y -= height
	} else if y < 0 {
		y += height
	}
	return pixels[y][x][0]
}

// genere des pixels avec une couleur random
func GenerateRandomPixels(grid_width int, grid_height int) [][][]uint8 {
	width = grid_width
	height = grid_height
	nestedPixels := make([][][]uint8, height)
	for y := range nestedPixels {
		nestedPixels[y] = make([][]uint8, width)
		for x := range nestedPixels[y] {
			aliveness := rand.Float32()
			nestedPixels[y][x] = []uint8{
				uint8(255 * aliveness), // R
				uint8(255 * aliveness), // G
				uint8(255 * aliveness), // B
			}
		}
	}

	return nestedPixels
}

func sigmoid(x float64, threshold float64, steepness float64) float64 {
	return 1 / (1 + math.Exp(-steepness*(x-threshold)))
}

func compute_new_state(S_n float64, S_m float64, B float64, S float64, K float64) uint8 {
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
    S_n := wheightedConvolve(pixels, r1, x, y)
    S_m := wheightedConvolve(pixels, r2, x, y)
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
