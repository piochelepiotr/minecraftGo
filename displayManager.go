package main
import (
	"fmt"
	"runtime"
    "log"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

type DisplayManager struct {
    window *glfw.Window
    windowWidth, windowHeight int
}

func (d *DisplayManager) createDisplay() {
    var err error
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
    d.window, err = glfw.CreateWindow(d.windowWidth, d.windowHeight, "Cube", nil, nil)
	if err != nil {
		panic(err)
	}
	d.window.MakeContextCurrent()

	// Initialize Glow
	if err := gl.Init(); err != nil {
		panic(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)
}

func (d *DisplayManager) updateDisplay() {
    gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
    d.window.SwapBuffers()
    glfw.PollEvents()
}

func (d *DisplayManager) closeDisplay() {
	glfw.Terminate()
}

