package ux

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/piochelepiotr/minecraftGo/game_engine/guis"
	"github.com/piochelepiotr/minecraftGo/game_engine/loader"
	"github.com/piochelepiotr/minecraftGo/game_engine/render"
)

const bottomBarHeight float32 = 0.1
const bottomBarItems = 9

type BottomBar struct {
	background      guis.GuiTexture
	items []guis.GuiTexture
}

func NewBottomBar(aspectRatio float32) *BottomBar {
	item := loader.LoadGuiTexture("textures/item.png", mgl32.Vec2{}, mgl32.Vec2{})
	items := make([]guis.GuiTexture, 0, bottomBarItems)
	for i := 0; i < bottomBarItems; i++ {
		items = append(items, guis.GuiTexture{Id: item.Id})
	}
	b := &BottomBar{
		background:      loader.LoadGuiTexture("textures/black.png", mgl32.Vec2{}, mgl32.Vec2{}),
		items: items,
	}
	b.Resize(aspectRatio)
	return b
}

func (b *BottomBar) Resize(aspectRatio float32) {
	itemWidth := bottomBarHeight / aspectRatio
	bottomBarWidth := float32(bottomBarItems) * itemWidth
	posY := 1-bottomBarHeight
	b.background.Scale = mgl32.Vec2{bottomBarWidth, bottomBarHeight}
	b.background.Position = mgl32.Vec2{0, posY}
	for i := range b.items {
		b.items[i].Scale = mgl32.Vec2{itemWidth, bottomBarHeight}
		b.items[i].Position = mgl32.Vec2{float32(i-bottomBarItems/2)*2 * itemWidth, posY}
	}
}

// Render renders all objects on the screen
func (b *BottomBar) Render(renderer *render.MasterRenderer) {
	renderer.ProcessGui(b.background)
	for _, item := range b.items {
		renderer.ProcessGui(item)
	}
}
