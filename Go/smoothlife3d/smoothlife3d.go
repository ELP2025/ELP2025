package smoothlife3d

import (
	"math"
	"math/rand"
  "sync"
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
func GenerateRandomPixels(grid_width, grid_height int, kernelRadius float64, threshold float32) ([]uint8, [][]float64, [][]float64, [][]float64) {
  ra = kernelRadius
  width = grid_width
	height = grid_height
  world1 := make([][]float64, height)
  world2 := make([][]float64, height)
  world3 := make([][]float64, height)
	nestedPixels := make([]uint8, height*width*3)
  for y := range world1 {
    world1[y] = make([]float64, width)
		world2[y] = make([]float64, width)
		world3[y] = make([]float64, width)
    for x := range world1[y] {
      if rand.Float32() < threshold {
        world1[y][x] = rand.Float64()
        world2[y][x] = rand.Float64()
        world3[y][x] = rand.Float64()
        nestedPixels[(y*height+x)*3] = uint8(255 * world1[y][x])
        nestedPixels[(y*height+x)*3+1] = uint8(255 * world2[y][x])
        nestedPixels[(y*height+x)*3+2] = uint8(255 * world3[y][x])
      } else {
        for c := 0; c < 3; c++ {
          nestedPixels[(y*height+x)*3+c] = uint8(0)
        } 
      }
		}
	}

	return nestedPixels, world1, world2, world3
}

func updateLine(world1, world2, world3 [][]float64, pixels []uint8, newWorld1, newWorld2, newWorld3 [][]float64, startLine, endLine int) {
  for y:= startLine; y<endLine; y++ {
    for x:=0; x < width; x++ {	
      outer := outerKernel(world1, x, y, int(ra-1))
		  inner := innerKernel(world1, x, y, int(ra-1)/3)
		  newWorld1[y][x] = clamp(world1[y][x] + dt*(2*s(outer, inner, b1, b2, d1, d2) - 1), 0, 1)
    
      outer = outerKernel(world2, x, y, int(ra-1))
		  inner = innerKernel(world2, x, y, int(ra-1)/3)
		  newWorld2[y][x] = clamp(world2[y][x] + dt*(2*s(outer, inner, b1, b2, d1, d2) - 1), 0, 1)
		
      outer = outerKernel(world3, x, y, int(ra-1))
		  inner = innerKernel(world3, x, y, int(ra-1)/3)
		  newWorld3[y][x] = clamp(world3[y][x] + dt*(2*s(outer, inner, b1, b2, d1, d2) - 1),0, 1)
  
      pixels[(y*height+x)*3] = uint8(255 * world1[y][x])
      pixels[(y*height+x)*3+1] = uint8(255 * world2[y][x])
      pixels[(y*height+x)*3+2] = uint8(255 * world3[y][x])
    }
  }
}

func UpdateGrid(pixels []uint8, world1, world2, world3 [][]float64) ([]uint8, [][]float64, [][]float64, [][]float64) {
	var wg sync.WaitGroup

  newWorld1 := make([][]float64, height)
  newWorld2 := make([][]float64, height)
  newWorld3 := make([][]float64, height)

  for i := 0; i < height; i++ {
		newWorld1[i] = make([]float64, width)
		newWorld2[i] = make([]float64, width)
		newWorld3[i] = make([]float64, width)
	}

  thread := 11
  linesPerThread := height/thread

	for i := 0; i < thread; i++ {
    startLine := i* linesPerThread
    endLine := startLine + linesPerThread
		
    if i == thread-1 {
      endLine = height
    }

    wg.Add(1)
		go func(startLine, endLine int) {
			defer wg.Done()
      updateLine(world1, world2, world3, pixels, newWorld1, newWorld2, newWorld3, startLine, endLine)
		}(startLine, endLine)
	}

	wg.Wait()

	return pixels, newWorld1, newWorld2, newWorld3
}
