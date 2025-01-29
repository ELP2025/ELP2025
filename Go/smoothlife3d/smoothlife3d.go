package smoothlife3d

import (
	"math"
	"math/rand"
  "sync"
  "runtime"
  "github.com/mjibson/go-dsp/fft"
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

  world1, world2, world3 [][]float64

  bigKernelFFT [][]complex128
  smallKernelFFT [][]complex128
)

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

func fftConvolve(worldFFT , kernelFFT [][]complex128) [][]float64 {

	resultFFT := make([][]complex128, height)
	for y := 0; y < height; y++ {
		resultFFT[y] = make([]complex128, width)
		for x := 0; x < width; x++ {
			resultFFT[y][x] = worldFFT[y][x] * kernelFFT[y][x]
		}
	}

	result := fft.IFFT2(resultFFT)
	return complexToReal(result)
}

func complexToReal(input [][]complex128) [][]float64 {
	output := make([][]float64, height)
	for y := 0; y < height; y++ {
		output[y] = make([]float64, width)
		for x := 0; x < width; x++ {
			output[y][x] = real(input[y][x])
		}
	}
	return output
}

func innerKernel(worldFFT [][]complex128) [][]float64 {
	return fftConvolve(worldFFT, smallKernelFFT)
}

func outerKernel(worldFFT [][]complex128) [][]float64 {
	return fftConvolve(worldFFT, bigKernelFFT)
}

func generateKernelFFT(radius float64) [][]complex128 {
	kernel := make([][]float64, height)
	centerX, centerY := width/2, height/2

	for y := 0; y < height; y++ {
		kernel[y] = make([]float64, width)
		for x := 0; x < width; x++ {
			dist := math.Sqrt(float64((x-centerX)*(x-centerX) + (y-centerY)*(y-centerY)))
			if dist <= radius {
				kernel[y][x] = math.Exp(-0.5 * (dist * dist) / (radius * radius))
			} else {
				kernel[y][x] = 0.0
			}
		}
	}

  kernelFFT := fft.FFT2Real(kernel)
	return kernelFFT
}

// genere des pixels avec une couleur random
func GenerateRandomPixels(grid_width, grid_height int, kernelRadius float64, threshold float32) []uint8 {
  ra = kernelRadius
  width = grid_width
	height = grid_height

  world1 = make([][]float64, height)
  world2 = make([][]float64, height)
  world3 = make([][]float64, height)
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
  bigKernelFFT = generateKernelFFT(ra)
  smallKernelFFT = generateKernelFFT(ra/3)
	return nestedPixels
}

func UpdateGrid(pixels []uint8) []uint8 {
	var wg sync.WaitGroup

  newWorld1 := world1
  newWorld2 := world2 
  newWorld3 := world3

  worldFFT1 := fft.FFT2Real(world1)
  worldFFT2 := fft.FFT2Real(world2)
  worldFFT3 := fft.FFT2Real(world3)

  outer1 := outerKernel(worldFFT1) 
  outer2 := outerKernel(worldFFT2)
  outer3 := outerKernel(worldFFT3)

  inner1 := innerKernel(worldFFT1)
  inner2 := innerKernel(worldFFT2)
  inner3 := innerKernel(worldFFT3)

  thread := runtime.NumCPU()*2
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
        for y:= startLine; y<endLine; y++ {
          for x:=0; x < width; x++ {	
		        newWorld1[y][x] = clamp(world1[y][x] + dt*(2*s(outer1[y][x], inner1[y][x], b1, b2, d1, d2) - 1), 0, 1)
		        newWorld2[y][x] = clamp(world2[y][x] + dt*(2*s(outer2[y][x], inner2[y][x], b1, b2, d1, d2) - 1), 0, 1)
		        newWorld3[y][x] = clamp(world3[y][x] + dt*(2*s(outer3[y][x], inner3[y][x], b1, b2, d1, d2) - 1),0, 1)
  
            pixels[(y*height+x)*3] = uint8(255 * newWorld1[y][x])
            pixels[(y*height+x)*3+1] = uint8(255 * newWorld2[y][x])
            pixels[(y*height+x)*3+2] = uint8(255 * newWorld3[y][x])
    }
  }

  world1 = newWorld1
  world2 = newWorld2
  world3 = newWorld3
		}(startLine, endLine)
	}

	wg.Wait()

	return pixels
}
