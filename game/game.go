package game

import (
	"fmt"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/piochelepiotr/minecraftGo/game_engine/loader"
	"github.com/piochelepiotr/minecraftGo/game_engine/render"
	"github.com/piochelepiotr/minecraftGo/state"
	"time"
)

// Game is the state in which the player is playing
type Game struct {
	loader *loader.Loader
	changeState chan state.Switch
	display     *render.DisplayManager
	state state.ID
	gamingState *GamingState
	inGameMenuState *InGameMenuState
	mainMenuState *MainMenuState
	inventoryState *InventoryState
	renderer *render.MasterRenderer

	structEditing bool
}

// Start starts the main event loop of the game
func Start(display *render.DisplayManager, structEditing bool) {
	changeState :=  make(chan state.Switch, 1)
	loader := loader.NewLoader()
	gameState := &Game{
		loader: loader,
		inGameMenuState: NewInGameMenuState(display, loader, changeState),
		changeState: changeState,
		display:     display,
		state: state.Empty,
		renderer: render.CreateMasterRenderer(display.AspectRatio(), loader),
		structEditing: structEditing,
	}

	display.ResizeCallBack = gameState.Resize

	gameState.switchState(state.Switch{ID: state.MainMenu})
	// gameState.switchState(state.Switch{ID: state.Game, WorldName: "World_0"})

	gameState.run()
	// switch to empty state to close correctly everything
	gameState.switchState(state.Switch{ID: state.Empty})
}

func (g *Game) Resize(aspectRatio float32) {
	g.renderer.Resize(aspectRatio)
	if g.state == state.GameMenu || g.state == state.Game {
		g.gamingState.Resize(aspectRatio)
	}
}

func (g *Game) run() {
	defer g.renderer.CleanUp()
	defer g.loader.CleanUp()

	updateTicker := time.NewTicker(time.Second)
	defer updateTicker.Stop()

	frames := 0
	for !g.display.Window.ShouldClose() {
		select {
		case <-updateTicker.C:
			g.Update()
			fmt.Printf("FPS is %d\n", frames)
			frames = 0
		case stateID := <-g.changeState:
			g.switchState(stateID)
		default:
			frames++
			g.Render(g.renderer)
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
			g.inventoryState = nil
		}
	default:
	}
	switch newState.ID {
	case state.Game:
		if !g.isInGame() {
			g.gamingState = NewGamingState(newState.WorldName, g.display, g.changeState, g.loader, g.structEditing)
			g.inventoryState = NewInventoryState(g.display, g.loader, g.gamingState.inventory, g.changeState)
		}
		g.display.Window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
		g.display.Window.SetKeyCallback(g.gamingState.keyCallback)
		g.display.Window.SetCursorPosCallback(g.gamingState.mouseMoveCallback)
		g.display.Window.SetMouseButtonCallback(g.gamingState.clickCallback)
		g.display.Window.SetScrollCallback(g.gamingState.scrollCallBack)
	case state.GameMenu:
		g.display.Window.SetInputMode(glfw.CursorMode, glfw.CursorNormal)
		g.display.Window.SetKeyCallback(g.inGameMenuState.keyCallback)
		g.display.Window.SetCursorPosCallback(g.inGameMenuState.mouseMoveCallback)
		g.display.Window.SetMouseButtonCallback(g.inGameMenuState.clickCallback)
		g.display.Window.SetScrollCallback(g.inGameMenuState.scrollCallBack)
	case state.MainMenu:
		g.mainMenuState = NewMainMenuState(g.display, g.loader, g.changeState)
		g.display.Window.SetInputMode(glfw.CursorMode, glfw.CursorNormal)
		g.display.Window.SetKeyCallback(g.mainMenuState.keyCallback)
		g.display.Window.SetCursorPosCallback(g.mainMenuState.mouseMoveCallback)
		g.display.Window.SetMouseButtonCallback(g.mainMenuState.clickCallback)
		g.display.Window.SetScrollCallback(g.mainMenuState.scrollCallBack)
	case state.Inventory:
		g.inventoryState.inventory.ReBuild()
		g.display.Window.SetInputMode(glfw.CursorMode, glfw.CursorNormal)
		g.display.Window.SetKeyCallback(g.inventoryState.keyCallback)
		g.display.Window.SetCursorPosCallback(g.inventoryState.mouseMoveCallback)
		g.display.Window.SetMouseButtonCallback(g.inventoryState.clickCallback)
		g.display.Window.SetScrollCallback(g.inventoryState.scrollCallBack)
	default:
	}
	g.state = newState.ID
}

func (g *Game) isInGame() bool {
	return g.state == state.GameMenu || g.state == state.Game || g.state == state.Inventory
}

// Update is called every second
func (g *Game) Update() {
	if g.isInGame() {
		g.gamingState.Update()
	}
}

// NextFrame makes time pass to move to the next frame of the game
func (g *Game) NextFrame() {
	if g.isInGame() {
		g.gamingState.NextFrame()
	}
}

// Render renders all objects on the screen
func (g *Game) Render(renderer *render.MasterRenderer) {
	if g.isInGame() {
		g.gamingState.Render(renderer)
	}
	if g.state == state.GameMenu {
		g.inGameMenuState.Render(renderer)
	}
	if g.state == state.MainMenu {
		g.mainMenuState.Render(renderer)
	}
	if g.state == state.Inventory {
		g.inventoryState.Render(renderer)
	}
	renderer.Render()
}