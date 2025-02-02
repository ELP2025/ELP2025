package main

import (
	"flag"
	"fmt"
	"image"
  _ "image/png"
  _ "image/jpeg"
	"log"
	"os"
	"runtime"
	"time"
	
	"main/opengl_utils"
	"main/smoothlife3d"

	"github.com/go-gl/glfw/v3.2/glfw"
)

// isPowerOfTwo returns true if n is a power of two.
func isPowerOfTwo(n int) bool {
	return n > 0 && (n&(n-1)) == 0
}

// loadImage loads an image from path and returns a []uint8 pixel slice in R,G,B format,
// its width, and height. It supports PNG and JPEG.
func loadImage(path string) ([]uint8, int, int, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, 0, 0, err
	}
	defer file.Close()

	// Decode the image (automatically detects format)
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, 0, 0, err
	}

	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()

	// Check if dimensions are powers of two.
	if !isPowerOfTwo(w) || !isPowerOfTwo(h) {
		return nil, 0, 0, fmt.Errorf("image dimensions (%d x %d) are not powers of two", w, h)
	}

	// Create a pixel slice in R,G,B format.
	pixels := make([]uint8, w*h*3)
	index := 0
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			// The RGBA() method returns values in the range [0, 65535].
			// Convert them to [0,255].
			pixels[index] = uint8(r >> 8)
			pixels[index+1] = uint8(g >> 8)
			pixels[index+2] = uint8(b >> 8)
			index += 3
		}
	}
	return pixels, w, h, nil
}

func main() {
	runtime.LockOSThread()

	// Define command-line flags.
	imagePath := flag.String("i", "", "path to image file to use as start grid")
	randomFlag := flag.Bool("r", false, "use random grid (requires -w and -h)")
	widthFlag := flag.Int("w", 1024, "grid width (must be power of two)")
	heightFlag := flag.Int("h", 1024, "grid height (must be power of two)")
  radiusFlag := flag.Float64("ra", 11, "radius to use for the outer kernel")
  thresholdFlag := flag.Float64("t", 1.00, "threshold for random grid generation")
  flag.Parse()

	var pixels []uint8
	var gridWidth, gridHeight int
	var kernelRadius = *radiusFlag
	var threshold = float32(*thresholdFlag)

	// Check that either an image or random mode is selected.
	if *imagePath != "" {
		// Load image and error-check dimensions.
		var err error
		pixels, gridWidth, gridHeight, err = loadImage(*imagePath)
		if err != nil {
			log.Fatalf("Error loading image: %v", err)
		}
		// Initialize the simulation state from the loaded image.
		pixels = smoothlife3d.LoadImagePixels(pixels, gridWidth, gridHeight, kernelRadius)
	} else if *randomFlag {
		// Check that the provided dimensions are powers of two.
		if !isPowerOfTwo(*widthFlag) || !isPowerOfTwo(*heightFlag) {
			log.Fatalf("Provided dimensions (%d x %d) are not powers of two", *widthFlag, *heightFlag)
		}
		gridWidth = *widthFlag
		gridHeight = *heightFlag
		pixels = smoothlife3d.GenerateRandomPixels(gridWidth, gridHeight, kernelRadius, threshold)
	} else {
		log.Fatalf("You must specify either an image (-i /path/to/image.png) or random mode (-r with -w (width) and -h (height), optionnaly -t (threshold value)). \n For both options, -ra specify the kernel radius")
	}

	// Initialize the OpenGL window with the chosen dimensions.
	window := opengl_utils.InitWindow(gridWidth, gridHeight)
	defer glfw.Terminate() // Ensure the window is closed properly

	for !window.ShouldClose() {
		t := time.Now()

		opengl_utils.UpdateTexture(pixels)
		pixels = smoothlife3d.UpdateGrid(pixels)

		fmt.Println("Last frame took", time.Since(t), "to render. Running at", 1.0/time.Since(t).Seconds(), "fps")
	}
}
