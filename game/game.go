package game

import (
	"fmt"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/piochelepiotr/minecraftGo/loader"
	"github.com/piochelepiotr/minecraftGo/render"
	"github.com/piochelepiotr/minecraftGo/state"
	"time"
)

// Game is the state in which the player is playing
type Game struct {
	changeState chan state.Switch
	display     *render.DisplayManager
	state state.ID
	gamingState *GamingState
	inGameMenuState *InGameMenuState
	mainMenuState *MainMenuState
}

// Start starts the main event loop of the game
func Start(display *render.DisplayManager) {
	changeState :=  make(chan state.Switch, 1)
	gameState := &Game{
		inGameMenuState: NewInGameMenuState(display, changeState),
		changeState: changeState,
		display:     display,
		state: state.Empty,
	}

	gameState.switchState(state.Switch{ID: state.MainMenu})

	gameState.run()
	// switch to empty state to close correctly everything
	gameState.switchState(state.Switch{ID: state.Empty})
}

func (g *Game) run() {
	renderer := render.CreateMasterRenderer()
	defer renderer.CleanUp()
	defer loader.CleanUp()

	updateTicker := time.NewTicker(time.Second)
	defer updateTicker.Stop()

	frames := 0
	for !g.display.Window.ShouldClose() {
		select {
		case <-updateTicker.C:
			start := time.Now()
			g.Update()
			stopTime := time.Now().Sub(start)
			fmt.Println(stopTime)
			fmt.Printf("FPS is %d\n", frames)
			frames = 0
		case stateID := <-g.changeState:
			g.switchState(stateID)
		default:
			frames++
			g.Render(renderer)
			g.NextFrame()
			g.display.UpdateDisplay()
			glfw.PollEvents()
		}
	}

}

func (g *Game) switchState(newState state.Switch) {
	switch g.state  {
	case state.MainMenu:
		g.mainMenuState = nil
	case state.Game:
		g.gamingState.pause()
	case state.GameMenu:
		if newState.ID != state.Game {
			g.gamingState.Close()
			g.gamingState = nil
		}
	default:
	}
	switch newState.ID {
	case state.Game:
		if g.state != state.GameMenu {
			g.gamingState = NewGamingState(newState.WorldName, g.display, g.changeState)
		}
		g.display.Window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
		g.display.Window.SetKeyCallback(g.gamingState.keyCallback)
		g.display.Window.SetCursorPosCallback(g.gamingState.mouseMoveCallback)
		g.display.Window.SetMouseButtonCallback(g.gamingState.clickCallback)
	case state.GameMenu:
		g.display.Window.SetInputMode(glfw.CursorMode, glfw.CursorNormal)
		g.display.Window.SetKeyCallback(g.inGameMenuState.keyCallback)
		g.display.Window.SetCursorPosCallback(g.inGameMenuState.mouseMoveCallback)
		g.display.Window.SetMouseButtonCallback(g.inGameMenuState.clickCallback)
	case state.MainMenu:
		g.mainMenuState = NewMainMenuState(g.display, g.changeState)
		g.display.Window.SetInputMode(glfw.CursorMode, glfw.CursorNormal)
		g.display.Window.SetKeyCallback(g.mainMenuState.keyCallback)
		g.display.Window.SetCursorPosCallback(g.mainMenuState.mouseMoveCallback)
		g.display.Window.SetMouseButtonCallback(g.mainMenuState.clickCallback)
	default:
	}
	g.state = newState.ID
}

// Update is called every second
func (g *Game) Update() {
	if g.state == state.GameMenu || g.state == state.Game {
		g.gamingState.Update()
	}
}

// NextFrame makes time pass to move to the next frame of the game
func (g *Game) NextFrame() {
	if g.state == state.GameMenu || g.state == state.Game {
		g.gamingState.NextFrame()
	}
}

// Render renders all objects on the screen
func (g *Game) Render(renderer *render.MasterRenderer) {
	if g.state == state.Game || g.state == state.GameMenu {
		g.gamingState.Render(renderer)
	}
	if g.state == state.GameMenu {
		g.inGameMenuState.Render(renderer)
	}
	if g.state == state.MainMenu {
		g.mainMenuState.Render(renderer)
	}
	renderer.Render()
}