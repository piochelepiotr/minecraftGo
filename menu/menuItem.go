package menu

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/piochelepiotr/minecraftGo/font"
	pfont "github.com/piochelepiotr/minecraftGo/font"
	"github.com/piochelepiotr/minecraftGo/guis"
	"github.com/piochelepiotr/minecraftGo/loader"
)

const (
	ItemHeight  float32 = 0.1
	MenuSpacing float32 = 0.05
	ItemWidth   float32 = 0.9
)

type MenuItem struct {
	text            font.GUIText
	guiTexture      guis.GuiTexture
	selectedTexture guis.GuiTexture
	index           int
}

func CreateMenuItem(text string, index int, font *pfont.FontType) *MenuItem {
	return &MenuItem{
		text:            loader.LoadText(pfont.CreateGUIText(text, 2, font, mgl32.Vec2{0, 0}, 1, true, ItemHeight, true)),
		index:           index,
		guiTexture:      loader.LoadGuiTexture("textures/stone.png", mgl32.Vec2{0, 0}, mgl32.Vec2{ItemWidth, ItemHeight}),
		selectedTexture: loader.LoadGuiTexture("textures/grass.png", mgl32.Vec2{0, 0}, mgl32.Vec2{ItemWidth, ItemHeight}),
	}
}

func getStartMenu(numberOfItems int) float32 {
	menuHeight := (float32(numberOfItems) - 1.0) * (ItemHeight + MenuSpacing)
	return -menuHeight/2 - ItemHeight/2
}

func blockSize() float32 {
	return ItemHeight + MenuSpacing
}

func itemIndex(x, y float32, numberOfItems int) int {
	if math.Abs(float64(x)) > float64(ItemWidth/2) {
		return -1
	}
	y = y - getStartMenu(numberOfItems)
	if y < 0 {
		return -1
	}
	index := int(y / blockSize())
	if index >= numberOfItems {
		return -1
	}
	if y-float32(index)*blockSize() > ItemHeight {
		return -1
	}
	return index
}

func (i *MenuItem) computeYPos(numberOfItems int) {
	yPos := getStartMenu(numberOfItems) + float32(i.index)*blockSize()
	yPos = 2 * yPos
	//yPos = -0.5
	i.guiTexture.Position = mgl32.Vec2{0, yPos + ItemHeight}
	i.selectedTexture.Position = mgl32.Vec2{0, yPos + ItemHeight}
	i.text.Position = mgl32.Vec2{0, 1 + yPos}
}
