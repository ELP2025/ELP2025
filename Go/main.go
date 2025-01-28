package main

import (
	"runtime"
	"time"
  "fmt"

	"main/opengl_utils"
	//"main/game_of_life"
	//"main/smoothlife"
  "main/smoothlife3d"

	"github.com/go-gl/glfw/v3.2/glfw"
)

const (
	width  = 800
	height = 800

	threshold = 0.80
)

func main() {
	runtime.LockOSThread()

	window := opengl_utils.InitWindow(width, height)
	defer glfw.Terminate() //Making sure we kill our window properly

	//nestedPixels := game_of_life.GenerateRandomPixels(width, height,threshold)
	pixels, world1, world2, world3 := smoothlife3d.GenerateRandomPixels(width, height, 11, threshold)

	for !window.ShouldClose() {
    t := time.Now()

		opengl_utils.UpdateTexture(pixels)
		// Dynamically update pixel data (optional)
		//nestedPixels = game_of_life.UpdateGrid(nestedPixels)
		pixels, world1, world2, world3 = smoothlife3d.UpdateGrid(pixels, world1, world2, world3)

    fmt.Println("Last frame took ", time.Since(t), "to render. Running at ", 1.0 / time.Since(t).Seconds(), "fps")
		//time.Sleep(time.Second/time.Duration(fps) - time.Since(t))
	}
}
