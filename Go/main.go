package main

import (
	"runtime"
	//"time"

	"main/opengl_utils"
	//"main/game_of_life"
	"main/smoothlife"

	"github.com/go-gl/glfw/v3.2/glfw"
)

const (
	width  = 200
	height = 200
	fps    = 60

	threshold = 1.0
)

func main() {
	runtime.LockOSThread()

	window := opengl_utils.InitWindow(width, height)
	defer glfw.Terminate() //Making sure we kill our window properly

	//nestedPixels := game_of_life.GenerateRandomPixels(width, height,threshold)
	pixels, world := smoothlife.GenerateRandomPixels(width, height, 11, threshold)

	for !window.ShouldClose() {
    //t := time.Now()

		opengl_utils.UpdateTexture(pixels)
		// Dynamically update pixel data (optional)
		//nestedPixels = game_of_life.UpdateGrid(nestedPixels)
		pixels, world = smoothlife.UpdateGrid(pixels, world)

		//time.Sleep(time.Second/time.Duration(fps) - time.Since(t))
	}
}
