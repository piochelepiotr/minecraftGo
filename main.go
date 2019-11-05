package main

import (
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/piochelepiotr/minecraftGo/game"
	"github.com/piochelepiotr/minecraftGo/render"
	"github.com/piochelepiotr/minecraftGo/state"
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"
)

const windowWidth = 800
const windowHeight = 600
const aspectRatio = windowWidth / windowHeight

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	d := render.DisplayManager{WindowWidth: windowWidth, WindowHeight: windowHeight}
	defer d.CloseDisplay()
	d.CreateDisplay()

	changeState := make(chan state.StateID, 1)
	defer close(changeState)

	gameState := game.NewGameState(aspectRatio, changeState)
	defer gameState.Close()

	resizeWindow := func(w *glfw.Window, width int, height int) {
		d.Resize(width, height)
	}

	clickCallback := func(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
		if action == glfw.Press {
			if button == glfw.MouseButtonRight {
				gameState.RightClick()
			} else if button == glfw.MouseButtonLeft {
				gameState.LeftClick()
			}
		}
	}

	mouseMoveCallback := func(w *glfw.Window, xpos float64, ypos float64) {
		x, y := d.GLPos(xpos, ypos)
		gameState.MouseMove(x, y)
	}

	keyCallback := func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		if action == glfw.Press {
			gameState.KeyPressed(key)
		} else if action == glfw.Release {
			gameState.KeyReleased(key)
		}
	}

	d.Window.SetKeyCallback(keyCallback)
	d.Window.SetCursorPosCallback(mouseMoveCallback)
	d.Window.SetMouseButtonCallback(clickCallback)
	d.Window.SetSizeCallback(resizeWindow)


	updateTicker := time.NewTicker(time.Second)
	defer updateTicker.Stop()

	for !d.Window.ShouldClose() {
		select {
			case <-updateTicker.C:
				gameState.Update()
			case stateID := <- changeState:
				if stateID == state.Game {
					d.Window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
					gameState.CloseMenu()
				} else if stateID == state.GameMenu {
					d.Window.SetInputMode(glfw.CursorMode, glfw.CursorNormal)
					gameState.OpenMenu()
				}
			default:
				gameState.NextFrame()
				d.UpdateDisplay()
				glfw.PollEvents()
		}
	}
}
