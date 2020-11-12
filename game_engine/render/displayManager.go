package render

import (
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
	Window         *glfw.Window
	width, height  int
	ResizeCallBack func(aspectRatio float32)
}

// NewDisplay create a glfw window
func NewDisplay(windowWidth, windowHeight int) *DisplayManager{
	d := DisplayManager{width: windowWidth, height: windowHeight}
	var err error
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}

	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	d.Window, err = glfw.CreateWindow(d.width, d.height, "Minecraft", nil, nil)
	if err != nil {
		panic(err)
	}
	d.Window.MakeContextCurrent()

	// Initialize Glow
	if err := gl.Init(); err != nil {
		panic(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Printf("OpenGL version: %s\n", version)
	d.Window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
	resizeWindow := func(w *glfw.Window, width int, height int) {
		d.Resize(width, height)
	}
	d.Window.SetSizeCallback(resizeWindow)
	return &d
}

func (d *DisplayManager) AspectRatio() float32 {
	return float32(d.width) / float32(d.height)
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
	xpos := float32(x) / float32(d.width)
	ypos := float32(y) / float32(d.height)
	xpos = xpos - 0.5
	ypos = ypos - 0.5
	return xpos, ypos
}

// Resize changes the size of display manager
func (d *DisplayManager) Resize(width, height int) {
	d.width = width
	d.height = height
	if d.ResizeCallBack != nil {
		d.ResizeCallBack(d.AspectRatio())
	}
}
