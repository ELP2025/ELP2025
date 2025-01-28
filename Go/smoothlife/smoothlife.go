package smoothlife

import (
	"math"
	"math/rand"
)

var (
  width = 1000
  height = 1000

  ra    float64 = 11
	alpha float64 = 0.028
	dt    float64 = 0.05

	b1 float64 = 0.278
	b2 float64 = 0.365
	d1 float64 = 0.267
	d2 float64 = 0.445
)

func modulo(a, b int) int {
  // Used to wrap values arround the grid
  return (a%b + b) % b
}

func clamp(x, min, max float64) float64 {
  // Make sure our values are not out of bound
  if x > max { 
    return max
  } else if x < min {
    return min
  }
  return x
}

func sigma1(x, a float64) float64 {
	return 1.0 / (1.0 + math.Exp(-(x-a)*4/alpha))
}

func sigma2(x, a, b float64) float64 {
	return sigma1(x, a) * (1 - sigma1(x, b))
}

func sigmam(x, y, m float64) float64 {
	return x*(1-sigma1(m, 0.5)) + y*sigma1(m, 0.5)
}

func s(n, m, b1, b2, d1, d2 float64) float64 {
	return sigma2(n, sigmam(b1, d1, m), sigmam(b2, d2, m))
}

func innerKernel(world [][]float64, x, y, radius int) float64{
  return convolve(world, x, y, radius, true)
}

func outerKernel(world [][]float64, x, y, radius int) float64 {
  return convolve(world, x, y, radius, false)
}

func convolve(world [][]float64, x, y, radius int, noCenter bool) float64 {
	sum := 0.0
  total := 0.0 
  for i := modulo(y, height) - radius; i <= modulo(y, height) + radius; i++ {
    for j:= modulo(x, width) - radius; j <= modulo(x, width) + radius; j++ {
      if !(noCenter && y == i && x == j) {
      dist := math.Sqrt(float64(i-y)*float64(i-y) + float64(j-x)*float64(j-x))
      wheight := math.Exp(-0.5 * math.Pow(dist/float64(radius), 2))

      sum += wheight * world[modulo(i, width)][modulo(j, height)]
      total += wheight
    }
    }
  }
  
  if total != 0 {
    return (sum / total)
  }
  return 0.0
}

// genere des pixels avec une couleur random
func GenerateRandomPixels(grid_width, grid_height, kernelRadius int, threshold float32) ([]uint8, [][]float64) {
  width = grid_width
	height = grid_height
  world := make([][]float64, height)
	nestedPixels := make([]uint8, height*width*3)
  for y := range world {
    world[y] = make([]float64, width)
		for x := range world[y] {
      if rand.Float32() < threshold {
        world[y][x] = rand.Float64()
        for c := 0; c < 3; c++ {
          nestedPixels[(y*height+x)*3+c] = uint8(255 * world[y][x])
        }
      } else {
        for c := 0; c < 3; c++ {
          nestedPixels[(y*height+x)*3+c] = uint8(0)
        } 
      }
		}
	}

	return nestedPixels, world
}

func updateLine(world [][]float64, y int) []float64 {
	newW := make([]float64, width)
	for x := range newW {
    outer := outerKernel(world, y, x,int(ra-1))
    inner := innerKernel(world, y, x,int(ra-1)/3)
    
    newW[x] = 2*s(outer, inner, b1, b2, d1, d2) - 1
	}
	return newW
}

func UpdateGrid(pixels []uint8, world [][]float64) ([]uint8, [][]float64) {
	newWorld := make([][]float64, height)
	done := make(chan int, height) // Channel to synchronize goroutines

	for i := 0; i < height; i++ {
		i := i // Capture the loop variable
		go func() {
			newWorld[i] = updateLine(world, i)
			done <- i // Signal completion
		}()
	}

	// Wait for all goroutines to finish
	for i := 0; i < height; i++ {
		<-done
	}

  for i:=range world {
    for j := range world[i] {
      world[i][j] += dt * newWorld[i][j]
      world[i][j] = clamp(world[i][j], 0, 1)
      for c :=0; c < 3; c++ {
        pixels[(i*height+j)*3+c] = uint8(255*world[i][j])
      }
    }
  }

	return pixels, world
}
