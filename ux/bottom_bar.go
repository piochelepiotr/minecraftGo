package ux

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/piochelepiotr/minecraftGo/game_engine/guis"
	"github.com/piochelepiotr/minecraftGo/game_engine/loader"
	"github.com/piochelepiotr/minecraftGo/game_engine/render"
)

const bottomBarHeight float32 = 0.1
const bottomBarWidth  float32 = 0.5

type BottomBar struct {
	background      guis.GuiTexture
}

func NewBottomBar() *BottomBar {
	posX := float32(0)
	posY := 1-bottomBarHeight
	return &BottomBar{
		background:      loader.LoadGuiTexture("textures/black.png", mgl32.Vec2{posX, posY}, mgl32.Vec2{bottomBarWidth, bottomBarHeight}),
	}
}

// Render renders all objects on the screen
func (b *BottomBar) Render(renderer *render.MasterRenderer) {
	renderer.ProcessGui(b.background)
}
