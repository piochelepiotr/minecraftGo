package menu

import (
	"github.com/piochelepiotr/minecraftGo/fontMeshCreator"
	"github.com/piochelepiotr/minecraftGo/guis"
	"github.com/piochelepiotr/minecraftGo/renderEngine"
)

type Menu struct {
	Opened    bool
	MenuItems []*MenuItem
	font      *fontMeshCreator.FontType
}

func CreateMenu(aspectRatio float32) Menu {
	return Menu{
		Opened:    false,
		MenuItems: make([]*MenuItem, 0),
		font:      renderEngine.LoadFont("./res/font.png", "./res/font.fnt", aspectRatio),
	}
}

func (m *Menu) AddItem(text string) {
	m.MenuItems = append(m.MenuItems, CreateMenuItem(text, len(m.MenuItems), m.font))
	for _, item := range m.MenuItems {
		item.computeYPos(len(m.MenuItems))
	}
}

func (m *Menu) GetMenuItems() []guis.GuiTexture {
	guis := make([]guis.GuiTexture, 0)
	for _, item := range m.MenuItems {
		guis = append(guis, item.guiTexture)
	}
	return guis
}

func (m *Menu) GetMenuTexts() []fontMeshCreator.GUIText {
	texts := make([]fontMeshCreator.GUIText, 0)
	for _, item := range m.MenuItems {
		texts = append(texts, item.text)
	}
	return texts
}
