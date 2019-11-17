package main

import (
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/piochelepiotr/minecraftGo/game"
	"github.com/piochelepiotr/minecraftGo/render"
	"github.com/piochelepiotr/minecraftGo/state"
	"log"
	"net/http"
	_ "net/http/pprof"
)

const windowWidth = 400
const windowHeight = 300
const aspectRatio = windowWidth / windowHeight

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	d := &render.DisplayManager{WindowWidth: windowWidth, WindowHeight: windowHeight}
	d.CreateDisplay()
	defer d.CloseDisplay()

	changeState := make(chan state.StateID, 1)
	defer close(changeState)

	resizeWindow := func(w *glfw.Window, width int, height int) {
		d.Resize(width, height)
	}

	d.Window.SetSizeCallback(resizeWindow)
	game.Run(aspectRatio, changeState, d)
}
