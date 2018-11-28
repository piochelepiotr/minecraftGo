package render

import (
	"fmt"
	"log"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

// DisplayManager manages the glfw window
type DisplayManager struct {
	Window                    *glfw.Window
	WindowWidth, WindowHeight int
}

// CreateDisplay create a glfw window
func (d *DisplayManager) CreateDisplay() {
	var err error
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}

	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	d.Window, err = glfw.CreateWindow(d.WindowWidth, d.WindowHeight, "Minecraft", nil, nil)
	if err != nil {
		panic(err)
	}
	d.Window.MakeContextCurrent()

	// Initialize Glow
	if err := gl.Init(); err != nil {
		panic(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)
	d.Window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
}

//UpdateDisplay polls events and swap buffers
func (d *DisplayManager) UpdateDisplay() {
	d.Window.SwapBuffers()
	glfw.PollEvents()
}

// CloseDisplay closes the glfw window
func (d *DisplayManager) CloseDisplay() {
	glfw.Terminate()
}

// GLPos returns opengl position of a point on the window
func (d *DisplayManager) GLPos(x, y float64) (float32, float32) {
	xpos := float32(x) / float32(d.WindowWidth)
	ypos := float32(y) / float32(d.WindowHeight)
	xpos = xpos - 0.5
	ypos = ypos - 0.5
	return xpos, ypos
}

// Resize changes the size of display manager
func (d *DisplayManager) Resize(width, height int) {
	d.WindowWidth = width
	d.WindowHeight = height
}
