package game

import (
	"fmt"
	"github.com/go-gl/glfw/v3.2/glfw"
	pmenu "github.com/piochelepiotr/minecraftGo/menu"
	"github.com/piochelepiotr/minecraftGo/render"
	"github.com/piochelepiotr/minecraftGo/state"
	"github.com/piochelepiotr/minecraftGo/world"
)
type MainMenuState struct {
	menu        *pmenu.Menu
	display     *render.DisplayManager
}

func loadWorlds() []string {
	return nil
}

func generateWorldName(worlds []string) string {
	contains := func(name string) bool {
		for _, n := range worlds {
			if n == name {
				return true
			}
		}
		return false
	}
	for i := 0;; i++ {
		worldName := fmt.Sprintf("World %d", i)
		if !contains(worldName) {
			return worldName
		}
	}
}

func NewMainMenuState(display *render.DisplayManager, changeState chan<- state.Switch) *MainMenuState{
	menu := pmenu.CreateMenu(display.AspectRatio())
	worlds := loadWorlds()
	for _, worldName := range worlds {
		menu.AddItem(worldName, func() { changeState <- state.Switch{ID: state.Game, WorldName: worldName} })
	}
	menu.AddItem("Create World", func() {
		name := generateWorldName(worlds)
		world.GenerateWorld(name)
		changeState <- state.Switch{ID: state.Game, WorldName: name}
	})
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