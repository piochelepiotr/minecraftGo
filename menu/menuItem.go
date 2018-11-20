package menu

import (
    "github.com/piochelepiotr/minecraftGo/guis"
    "github.com/piochelepiotr/minecraftGo/renderEngine"
	"github.com/go-gl/mathgl/mgl32"
)

const (
    ItemHeight float32 = 0.1
    MenuSpacing float32 = 0.05
    ItemWidth float32 = 0.9
)

type MenuItem struct {
    text string
    guiTexture guis.GuiTexture
    index int
}

func CreateMenuItem(text string, index int) *MenuItem {
    return &MenuItem{
        text: text,
        index: index,
        guiTexture: renderEngine.LoadGuiTexture("textures/stone.png", mgl32.Vec2{0, 0}, mgl32.Vec2{ItemWidth, ItemHeight}),
    }
}

func (i *MenuItem) computeYPos(numberOfItems int) {
    menuHeight := (float32(numberOfItems) - 1.0) * (ItemHeight*2) + (float32(numberOfItems) - 1.0) * (MenuSpacing*2)
    i.guiTexture.Position = mgl32.Vec2{0, -(menuHeight/2.0) + float32(i.index) * (ItemHeight + MenuSpacing) * 2}
}
