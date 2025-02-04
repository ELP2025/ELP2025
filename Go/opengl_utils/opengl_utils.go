package opengl_utils

import (
  "fmt"
  "strings"

  "github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

const (
  	vertexShaderSource = `
		#version 410
		layout(location = 0) in vec3 position; // Vertex position
		layout(location = 1) in vec2 texCoord; // Texture coordinates

		out vec2 fragTexCoord; // Pass texture coordinates to fragment shader

		void main() {
			fragTexCoord = texCoord;
			gl_Position = vec4(position, 1.0);
		}
	` + "\x00"

	fragmentShaderSource = `
		#version 410
		in vec2 fragTexCoord; // Texture coordinates from vertex shader
		out vec4 fragColor; // Output color

		uniform sampler2D screenTexture; // Texture containing pixel data

		void main() {
			fragColor = texture(screenTexture, fragTexCoord);
		}
	` + "\x00"
)

var (
	quadVertices = []float32{
		// Positions     // Texture Coords
		-1.0,  1.0, 0.0,  0.0, 1.0,
		-1.0, -1.0, 0.0,  0.0, 0.0,
		 1.0, -1.0, 0.0,  1.0, 0.0,

		-1.0,  1.0, 0.0,  0.0, 1.0,
		 1.0, -1.0, 0.0,  1.0, 0.0,
		 1.0,  1.0, 0.0,  1.0, 1.0,
	}

  width = 1000
  height = 1000
  
  window *glfw.Window
  vao uint32
  program uint32
  texture uint32

)

func InitWindow(window_width int, window_height int) *glfw.Window {
  // Init an OpenGL window with specific size
  width = window_width
  height = window_height
  window = initGlfw()

	program = initOpenGL()

	// Create a texture
	texture = createTexture()

	// Load quad vertices into VAO
	vao = makeVao(quadVertices)
  
  return window
}

func UpdateTexture(pixels []uint8) {
  // Edit the texture with pixels as argument
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexSubImage2D(gl.TEXTURE_2D, 0, 0, 0, int32(width), int32(height), gl.RGB, gl.UNSIGNED_BYTE, gl.Ptr(pixels))
  draw(window, program, vao, texture)
}

func draw(window *glfw.Window, program, vao, texture uint32) {
  // Draw the texture + shader to the screen
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(program)

	// Bind the texture
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)

	// Render the quad
	gl.BindVertexArray(vao)
	gl.DrawArrays(gl.TRIANGLES, 0, 6)

	glfw.PollEvents()
	window.SwapBuffers()
}

func initGlfw() *glfw.Window {
  // create a new GlFW instance to handle OpenGL
	if err := glfw.Init(); err != nil {
		panic(err)
	}
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(width, height, "smoothlife", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	return window
}

func initOpenGL() uint32 {
  // Init OpenGL vertexShader, fragmentShader and program
	if err := gl.Init(); err != nil {
		panic(err)
	}

	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}

	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	program := gl.CreateProgram()
	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	return program
}

func createTexture() uint32 {
  // creates a texture that will be on our quadVertices
	var texture uint32
	gl.GenTextures(1, &texture)
	gl.BindTexture(gl.TEXTURE_2D, texture)

	// Set texture parameters
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	// Initialize the texture with empty data
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGB, int32(width), int32(height), 0, gl.RGB, gl.UNSIGNED_BYTE, nil)

	return texture
}

func makeVao(vertices []float32) uint32 {
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	// Positions
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))

	// Texture coordinates
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))

	return vao
}

func compileShader(source string, shaderType uint32) (uint32, error) {
  // compile the shader defined in the var {} of this program. Just needed to init OpenGL
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}

