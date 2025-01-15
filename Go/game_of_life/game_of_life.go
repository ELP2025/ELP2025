package game_of_life

import (
  "math/rand"
)

var (
  width = 1000
  height = 1000
)


// genere des pixels avec une couleur random
func GenerateRandomPixels(grid_width int, grid_height int, threshold float32) [][][]uint8 {
	width = grid_width
  height = grid_height
  nestedPixels := make([][][]uint8, height)
	for y := range nestedPixels {
		nestedPixels[y] = make([][]uint8, width)
		for x := range nestedPixels[y] {
			if rand.Float32() < threshold { 
				nestedPixels[y][x] = []uint8{
					uint8(255), // R
					uint8(255), // G
					uint8(255), // B
				}
			} else {
				nestedPixels[y][x] = []uint8{0, 0, 0} // Default to black
      }
		}
	}
	return nestedPixels
}

func checkNeighbors(pixels [][][]uint8, x int, y int) int { 
    if x + 1 > width {
      x = 0
    } else if x-1 < 0 {
      x = width - 1
    }
    if y + 1 > height {
      y = 0
    } else if y-1 < 0 {
      y = height - 1
    }
    if pixels[y][x][0] == 255 {
      return 1 
    } else {
      return 0
    }
  }

func countNeighbors(pixels [][][]uint8, x int, y int) int {
  neighbors_count := 0 

  neighbors_count += checkNeighbors(pixels, x-1, y-1)
  neighbors_count += checkNeighbors(pixels, x, y-1)
  neighbors_count += checkNeighbors(pixels, x+1, y-1)
  neighbors_count += checkNeighbors(pixels, x-1, y)
  neighbors_count += checkNeighbors(pixels, x+1, y)
  neighbors_count += checkNeighbors(pixels, x-1, y+1)
  neighbors_count += checkNeighbors(pixels, x, y+1)
  neighbors_count += checkNeighbors(pixels, x+1, y+1)

  return neighbors_count
}

func updateLine(pixels [][][]uint8, y int) [][]uint8 {
  newPixels := make([][]uint8, width)
  
		for x := range newPixels {
      neigh := countNeighbors(pixels, x, y)
				if neigh < 2 || neigh > 3 {
          newPixels[x] = []uint8{0, 0, 0}
        } else if neigh == 2 {
          newPixels[x] = pixels[y][x]
        } else {
          newPixels[x] = []uint8{255, 255, 255}
        }
  }
  return newPixels
}

func UpdateGrid(pixels [][][]uint8) [][][]uint8{
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
