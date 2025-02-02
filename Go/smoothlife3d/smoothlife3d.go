package smoothlife3d

import (
	"math"
	"math/rand"
  "sync"
  "runtime"
  "main/fft" 
  "time"
  "fmt"
)

var (
  width = 1024
  height = 1024

	alpha float64 = 0.028
	dt    float64 = 0.15

	b1 float64 = 0.278
	b2 float64 = 0.365
	d1 float64 = 0.267
	d2 float64 = 0.445

  world1, world2, world3 []float64
  worldFFT1, worldFFT2, worldFFT3 []complex128

  bigKernelFFT []complex128
  smallKernelFFT []complex128

  kernelFunctions = []func() []float64{
    func() []float64 {return outerKernel(worldFFT1)},
    func() []float64 {return innerKernel(worldFFT1)},
    func() []float64 {return outerKernel(worldFFT2)},
    func() []float64 {return innerKernel(worldFFT2)},
    func() []float64 {return outerKernel(worldFFT3)},
    func() []float64 {return innerKernel(worldFFT3)},
  }
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

func fftConvolve(worldFFT , kernelFFT []complex128, threads int) []float64 {
  var wg sync.WaitGroup
  size := len(worldFFT)
  resultFFT := make([]complex128, size)

  indexPerThread := size/threads
  for t := 0; t < threads; t++ {
    wg.Add(1)
    start := t*indexPerThread
    end := (t+1)*indexPerThread
    if (t==threads-1) {end = size}

    go func(start, end int) {
      defer wg.Done()
      for i := 0; i < size; i++ {
		    resultFFT[i] = worldFFT[i] * kernelFFT[i]
	    }
    }(start, end)
  }

  wg.Wait()

	result := (fft.IFFT(resultFFT, height*width))
	return result
}

func innerKernel(worldFFT []complex128) []float64 {
	return fftConvolve(worldFFT, smallKernelFFT, 2)
}

func outerKernel(worldFFT []complex128) []float64 {
	return fftConvolve(worldFFT, bigKernelFFT, 2)
}

func generateKernelFFT(radius float64, skipCenter bool) []complex128 {
	kernel := make([]float64, height*width)
	centerX, centerY := width/2, height/2
  sum := 0.0 

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			dist := math.Sqrt(float64((x-centerX)*(x-centerX) + (y-centerY)*(y-centerY)))
			if dist <= radius {
				kernel[y*height+x] = math.Exp(-0.5 * (dist * dist) / (radius * radius))
			} else {
				kernel[y*height+x] = 0.0
			}
      sum += kernel[y*height+x]
		}  
	}

  if skipCenter {kernel[int(1/2 * width*(1+height))] = 0.0}

  if sum != 0 {
      for i := range kernel {
        kernel[i] /= sum
      }
    }

  kernelFFT := fft.FFT(kernel)
	return kernelFFT
}

// genere des pixels avec une couleur random
func GenerateRandomPixels(grid_width, grid_height int, kernelRadius float64, threshold float32) []uint8 {
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
  bigKernelFFT = generateKernelFFT(kernelRadius, false)
  smallKernelFFT = generateKernelFFT(kernelRadius/3, true)
	return nestedPixels
}

// LoadImagePixels initializes the simulation grid from an image’s pixel data.
// It assumes that pixels is a slice with length = gridWidth*gridHeight*3 (R,G,B),
// and that gridWidth and gridHeight are powers of two.
func LoadImagePixels(pixels []uint8, gridWidth, gridHeight int, kernelRadius float64) []uint8 {
	width = gridWidth
	height = gridHeight

	// Initialize the simulation world arrays.
	world1 = make([]float64, gridWidth*gridHeight)
	world2 = make([]float64, gridWidth*gridHeight)
	world3 = make([]float64, gridWidth*gridHeight)
	nestedPixels := make([]uint8, gridWidth*gridHeight*3)

	for y := 0; y < gridHeight; y++ {
		for x := 0; x < gridWidth; x++ {
			index := y*gridWidth + x
			// Copy the image pixel and normalize to [0,1] for the simulation state.
			nestedPixels[index*3] = pixels[index*3]
			nestedPixels[index*3+1] = pixels[index*3+1]
			nestedPixels[index*3+2] = pixels[index*3+2]
			world1[index] = float64(pixels[index*3]) / 255.0
			world2[index] = float64(pixels[index*3+1]) / 255.0
			world3[index] = float64(pixels[index*3+2]) / 255.0
		}
	}

	// Initialize the kernel FFTs as in the random generator.
	bigKernelFFT = generateKernelFFT(kernelRadius, false)
	smallKernelFFT = generateKernelFFT(kernelRadius/3, true)
	return nestedPixels
}

func UpdateGrid(pixels []uint8) []uint8 {
	var wg sync.WaitGroup

  t := time.Now()
  newWorld1 := make([]float64, len(world1))
  newWorld2 := make([]float64, len(world1))
  newWorld3 := make([]float64, len(world1))

  worldFFT1 = fft.FFT(world1)
  worldFFT2 = fft.FFT(world2)
  worldFFT3 = fft.FFT(world3)
  fmt.Println("Precomputation took ", time.Since(t))

  t = time.Now()
  convolutions := make([][]float64, 6)

  for i, fun := range kernelFunctions {
    wg.Add(1)
    go func(index int, fun func() []float64) {
      defer wg.Done()
      convolutions[index] = fun()
    }(i, fun)
  }

  wg.Wait()
  fmt.Println("Convolutions took ", time.Since(t))
 
  thread := runtime.NumCPU()*2
  linesPerThread := height/thread

  t = time.Now()
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
		        newWorld1[index] = clamp(world1[index] + dt*(2*s(convolutions[0][index], convolutions[1][index], b1, b2, d1, d2) - 1), 0, 1)
		        newWorld2[index] = clamp(world2[index] + dt*(2*s(convolutions[2][index], convolutions[3][index], b1, b2, d1, d2) - 1), 0, 1)
		        newWorld3[index] = clamp(world3[index] + dt*(2*s(convolutions[4][index], convolutions[5][index], b1, b2, d1, d2) - 1),0, 1)
  
            pixels[index*3] = uint8(255 * newWorld1[index])
            pixels[index*3+1] = uint8(255 * newWorld2[index])
            pixels[index*3+2] = uint8(255 * newWorld3[index])
    }
  }

		}(startLine, endLine)
	}
	wg.Wait()
  fmt.Println("Calculation new state took ", time.Since(t))

  t = time.Now()
  copy(world1, newWorld1)
  copy(world2, newWorld2)
  copy(world3, newWorld3)
  fmt.Println("Copying took ", time.Since(t))

	return pixels
}
