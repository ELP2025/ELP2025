package smoothlife

import (
	"math"
	"math/rand"
	"math"
	"math/rand"
)

var (
	width  = 1000
	height = 1000
var (
	width  = 1000
	height = 1000
)

func circleKernel(radius float32,center_x int, center_y int) []struct {Y; X int} {
  result := make([]struct{Y; X int})
  for x <= center_x - radius; x >= center_x + radius; x++ { //Reducing the interval for our search
    for y <= center_y - radius; y >= center_y + radius; y++ { //Same
      if sqrt((x-center_x)**2 + (y-center_y)**2) < radius { // If (x,y) is in circle of radius 
        result = append(result, struct {Y; X int}{Y : y, X : x})
      }
    }
  }
  return result
}

func wheightedConvolve(pixels [][][][]uint8, radius float32, x int, int y) float32 {
  points := circleKernel(radius, x, y)
  sum := 0.0
  count :=0 
  for _,point := range points {
    px, py = point.X, point.Y
    sum += getNeighborValue(pixels, x, y)
    count +=1
  }

  if count : return sum/count
  return 0

}

func getNeighborValue(pixels [][][]uint8, x int, y int) int {
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
}func compute_new_state(S_n float32, S_m float32, B float32, S float32, K float32) uint8 {
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
