package main

import (
	"runtime"
	"time"
  "fmt"
	"main/opengl_utils"
  "main/smoothlife3d"

	"github.com/go-gl/glfw/v3.2/glfw"
)

const (
	width  = 1024
	height = 1024

	threshold = 1.00
)

func main() {
	runtime.LockOSThread()

	window := opengl_utils.InitWindow(width, height)
	defer glfw.Terminate() //Making sure we kill our window properly

	pixels:= smoothlife3d.GenerateRandomPixels(width, height, 11, threshold)

	for !window.ShouldClose() {
    t := time.Now()

		opengl_utils.UpdateTexture(pixels)
		pixels = smoothlife3d.UpdateGrid(pixels)

    fmt.Println("Last frame took ", time.Since(t), "to render. Running at ", 1.0 / time.Since(t).Seconds(), "fps")
	}
}
