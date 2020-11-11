package game

import (
	"github.com/go-gl/glfw/v3.2/glfw"
	pmenu "github.com/piochelepiotr/minecraftGo/menu"
	"github.com/piochelepiotr/minecraftGo/render"
	"github.com/piochelepiotr/minecraftGo/state"
)
type MainMenuState struct {
	menu        *pmenu.Menu
	display     *render.DisplayManager
}

func NewMainMenuState(display *render.DisplayManager, changeState chan<- state.Switch) *MainMenuState{
	menu := pmenu.CreateMenu(display.AspectRatio())
	menu.AddItem("World 1", func() { changeState <- state.Switch{ID: state.Game} })
	menu.AddItem("World 2", func() { changeState <- state.Switch{ID: state.Game} })
	menu.AddItem("Create World", func() { changeState <- state.Switch{ID: state.Game} })
	menu.AddItem("Exit Game", func() { display.Window.SetShouldClose(true) })
	return &MainMenuState{
		menu: menu,
		display: display,
	}
}

func (s *MainMenuState) clickCallback(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
	if action == glfw.Press {
		if button == glfw.MouseButtonLeft {
			s.menu.LeftClick()
		}
	}
}

func (s *MainMenuState) mouseMoveCallback(w *glfw.Window, xpos float64, ypos float64) {
	x, y := s.display.GLPos(xpos, ypos)
	s.menu.ComputeSelectedItem(x, y)
}

func (s *MainMenuState) keyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
}

func (s *MainMenuState) Render(renderer *render.MasterRenderer) {
	s.menu.Render(renderer)
}