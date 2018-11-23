package renderEngine

import (
	"fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"log"
	"runtime"
)

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

type DisplayManager struct {
	Window                    *glfw.Window
	WindowWidth, WindowHeight int
}

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
}

func (d *DisplayManager) UpdateDisplay() {
	d.Window.SwapBuffers()
	glfw.PollEvents()
}

func (d *DisplayManager) CloseDisplay() {
	glfw.Terminate()
}
