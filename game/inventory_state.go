package game

import (
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/piochelepiotr/minecraftGo/game_engine/loader"
	"github.com/piochelepiotr/minecraftGo/game_engine/render"
	"github.com/piochelepiotr/minecraftGo/state"
	"github.com/piochelepiotr/minecraftGo/ux"
	"github.com/piochelepiotr/minecraftGo/world"
)
type InventoryState struct {
	display     *render.DisplayManager
	inventory *ux.Inventory
	exit func()
}

func NewInventoryState(display *render.DisplayManager, loader *loader.Loader, inventory *world.Inventory, changeState chan<- state.Switch) *InventoryState{
	return &InventoryState{
		display: display,
		inventory: ux.NewInventory(display.AspectRatio(), loader, inventory),
		exit: func() {
			changeState <- state.Switch{ID: state.Game}
		},
	}
}

func (s *InventoryState) clickCallback(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
	if action == glfw.Press {
		if button == glfw.MouseButtonLeft {
		}
	}
}

func (s *InventoryState) mouseMoveCallback(w *glfw.Window, xpos float64, ypos float64) {
	// x, y := s.display.GLPos(xpos, ypos)
}

func (s *InventoryState) keyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if action == glfw.Press {
		if key == glfw.KeyEscape {
			s.exit()
		}
	}
}

func (s *InventoryState) Render(renderer *render.MasterRenderer) {
	s.inventory.Render(renderer)
}

func (s *InventoryState) scrollCallBack(w *glfw.Window, xoff float64, yoff float64) {
}
