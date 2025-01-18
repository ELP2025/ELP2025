package main

import (
	"runtime"
	"time"

	"main/opengl_utils"
	//"main/game_of_life"
	"main/smoothlife"

	"github.com/go-gl/glfw/v3.2/glfw"
)

const (
	width  = 1000
	height = 1000
	fps    = 60

	threshold = 0.05
)

func main() {
	runtime.LockOSThread()

	window := opengl_utils.InitWindow(width, height)
	defer glfw.Terminate() //Making sure we kill our window properly

	//nestedPixels := game_of_life.GenerateRandomPixels(width, height,threshold)
	nestedPixels := smoothlife.GenerateRandomPixels(width, height, 4, 16, threshold)

	for !window.ShouldClose() {
		t := time.Now()

		opengl_utils.UpdateTexture(nestedPixels)
		// Dynamically update pixel data (optional)
		//nestedPixels = game_of_life.UpdateGrid(nestedPixels)
		nestedPixels = smoothlife.UpdateGrid(nestedPixels)

		time.Sleep(time.Second/time.Duration(fps) - time.Since(t))
	}
}
