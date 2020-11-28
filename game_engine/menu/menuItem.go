package menu

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/piochelepiotr/minecraftGo/game_engine/font"
	pfont "github.com/piochelepiotr/minecraftGo/game_engine/font"
	"github.com/piochelepiotr/minecraftGo/game_engine/guis"
	"github.com/piochelepiotr/minecraftGo/game_engine/loader"
)

const (
	//ItemHeight is the height of one menu item
	ItemHeight float32 = 0.1
	//MenuSpacing is the spacing between two items
	MenuSpacing float32 = 0.05
	//ItemWidth is the width of a menu item
	ItemWidth float32 = 0.9
)

// Item is an item in the menu
type Item struct {
	text            font.GUIText
	guiTexture      guis.GuiTexture
	selectedTexture guis.GuiTexture
	index           int
	callback        func()
}

// CreateItem creates text and gui for menu item
func CreateItem(text string, index int, font *pfont.FontType, callback func(), loader *loader.Loader) *Item {
	return &Item{
		text:            loader.LoadText(pfont.CreateGUIText(text, 2, font, mgl32.Vec2{0, 0}, 1, ItemHeight, true, mgl32.Vec3{})),
		index:           index,
		guiTexture:      loader.LoadGuiTexture("textures/stone.png", mgl32.Vec2{0, 0}, mgl32.Vec2{ItemWidth, ItemHeight}),
		selectedTexture: loader.LoadGuiTexture("textures/dark_stone.png", mgl32.Vec2{0, 0}, mgl32.Vec2{ItemWidth, ItemHeight}),
		callback: callback,
	}
}

func startY(numberOfItems int) float32 {
	menuHeight := float32(numberOfItems) * ItemHeight + float32(numberOfItems-1) * MenuSpacing
	return -pos(menuHeight/2)
}

func blockSize() float32 {
	return ItemHeight + MenuSpacing
}

func itemIndex(x, y float32, numberOfItems int) int {
	if math.Abs(float64(x)) > float64(ItemWidth/2) {
		return -1
	}
	y = pos(y)
	y = y - startY(numberOfItems)
	if y < 0 {
		return -1
	}
	index := int(y / pos(blockSize()))
	if index >= numberOfItems {
		return -1
	}
	if y-pos(float32(index)*blockSize()) > pos(ItemHeight) {
		return -1
	}
	return index
}

func pos(x float32) float32 {
	return x*2
}

func (i *Item) computeYPos(numberOfItems int) {
	yPos := startY(numberOfItems) + pos(float32(i.index)*(ItemHeight+MenuSpacing)) + pos(ItemHeight)/2
	i.guiTexture.Position = mgl32.Vec2{0, yPos}
	i.selectedTexture.Position = mgl32.Vec2{0, yPos}
	i.text.Position = mgl32.Vec2{-pos(i.text.Width/2), yPos-pos(i.text.GetLineHeight()/2)}
}
