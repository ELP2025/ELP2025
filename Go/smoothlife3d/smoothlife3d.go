package smoothlife3d

import (
	"math"
	"math/rand"
  "sync"
  "runtime"
  "github.com/mjibson/go-dsp/fft"

  "fmt"
  "time"
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

  world1, world2, world3 []float64

  bigKernelFFT []complex128
  smallKernelFFT []complex128
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

func fftConvolve(worldFFT , kernelFFT []complex128) []float64 {

	resultFFT := make([]complex128, height*width)
	for i := 0; i < height; i++ {
			resultFFT[i] = worldFFT[i] * kernelFFT[i]
	}
	

	result := fft.IFFT(resultFFT)
	return complexToReal(result)
}

func complexToReal(input []complex128) []float64 {
  t := time.Now()
	output := make([]float64, height*width)
	for i := 0; i < height*width; i++ {
			output[i] = real(input[i])
	}
  fmt.Println("ComplexToReal took : ", time.Since(t))
	return output
}

func innerKernel(worldFFT []complex128) []float64 {
	return fftConvolve(worldFFT, smallKernelFFT)
}

func outerKernel(worldFFT []complex128) []float64 {
	return fftConvolve(worldFFT, bigKernelFFT)
}

func generateKernelFFT(radius float64) []complex128 {
	kernel := make([]float64, height*width)
	centerX, centerY := width/2, height/2

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			dist := math.Sqrt(float64((x-centerX)*(x-centerX) + (y-centerY)*(y-centerY)))
			if dist <= radius {
				kernel[y*height+x] = math.Exp(-0.5 * (dist * dist) / (radius * radius))
			} else {
				kernel[y*height+x] = 0.0
			}
		}
	}

  kernelFFT := fft.FFTReal(kernel)
	return kernelFFT
}

// genere des pixels avec une couleur random
func GenerateRandomPixels(grid_width, grid_height int, kernelRadius float64, threshold float32) []uint8 {
  fft.SetWorkerPoolSize(0)
  ra = kernelRadius
  width = grid_width
	height = grid_height

  world1 = make([]float64, height*width)
  world2 = make([]float64, height*width)
  world3 = make([]float64, height*width)
	nestedPixels := make([]uint8, height*width*3)

  for y := range height {

    for x := range width {
      if rand.Float32() < threshold {
        index := y*height+x
        world1[index] = rand.Float64()
        world2[index] = rand.Float64()
        world3[index] = rand.Float64()

        nestedPixels[index*3] = uint8(255 * world1[index])
        nestedPixels[index*3+1] = uint8(255 * world2[index])
        nestedPixels[index*3+2] = uint8(255 * world3[index])
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

  t := time.Now()
  newWorld1 := world1
  newWorld2 := world2 
  newWorld3 := world3
  fmt.Println("Copying worlds took : ", time.Since(t))
  t = time.Now()

  worldFFT1 := fft.FFTReal(world1)
  worldFFT2 := fft.FFTReal(world2)
  worldFFT3 := fft.FFTReal(world3)
  fmt.Println("Computing FFT took : ", time.Since(t))
  t= time.Now()

  outer1 := outerKernel(worldFFT1) 
  outer2 := outerKernel(worldFFT2)
  outer3 := outerKernel(worldFFT3)
  fmt.Println("Outer kernels took : ", time.Since(t))
  t = time.Now()

  inner1 := innerKernel(worldFFT1)
  inner2 := innerKernel(worldFFT2)
  inner3 := innerKernel(worldFFT3)
  fmt.Println("Inner Kernels took : ", time.Since(t))

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
            index := y*height+x
		        newWorld1[index] = clamp(world1[index] + dt*(2*s(outer1[index], inner1[index], b1, b2, d1, d2) - 1), 0, 1)
		        newWorld2[index] = clamp(world2[index] + dt*(2*s(outer2[index], inner2[index], b1, b2, d1, d2) - 1), 0, 1)
		        newWorld3[index] = clamp(world3[index] + dt*(2*s(outer3[index], inner3[index], b1, b2, d1, d2) - 1),0, 1)
  
            pixels[index*3] = uint8(255 * newWorld1[index])
            pixels[index*3+1] = uint8(255 * newWorld2[index])
            pixels[index*3+2] = uint8(255 * newWorld3[index])
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
