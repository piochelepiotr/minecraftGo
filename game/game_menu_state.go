package game

import (
	"github.com/go-gl/glfw/v3.2/glfw"
	pmenu "github.com/piochelepiotr/minecraftGo/menu"
	"github.com/piochelepiotr/minecraftGo/render"
	"github.com/piochelepiotr/minecraftGo/state"
)
type InGameMenuState struct {
	menu        *pmenu.Menu
	display     *render.DisplayManager
}

func NewInGameMenuState(display *render.DisplayManager, changeState chan<- state.Switch) *InGameMenuState{
	menu := pmenu.CreateMenu(display.AspectRatio())
	menu.AddItem("Resume game", func() { changeState <- state.Switch{ID: state.Game} })
	menu.AddItem("Exit game", func() { display.Window.SetShouldClose(true) })
	menu.AddItem("Go to main menu", func() { changeState <- state.Switch{ID: state.MainMenu} })
	return &InGameMenuState{
		menu: menu,
		display: display,
	}
}

func (s *InGameMenuState) clickCallback(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
	if action == glfw.Press {
		if button == glfw.MouseButtonLeft {
			s.menu.LeftClick()
		}
	}
}

func (s *InGameMenuState) mouseMoveCallback(w *glfw.Window, xpos float64, ypos float64) {
	x, y := s.display.GLPos(xpos, ypos)
	s.menu.ComputeSelectedItem(x, y)
}

func (s *InGameMenuState) keyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
}

func (s *InGameMenuState) Render(renderer *render.MasterRenderer) {
	s.menu.Render(renderer)
}