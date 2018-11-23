package menu

import (
	"github.com/piochelepiotr/minecraftGo/fontMeshCreator"
	"github.com/piochelepiotr/minecraftGo/guis"
	"github.com/piochelepiotr/minecraftGo/loader"
)

type Menu struct {
	Opened       bool
	MenuItems    []*MenuItem
	font         *fontMeshCreator.FontType
	SelectedItem int
}

func CreateMenu(aspectRatio float32) Menu {
	return Menu{
		Opened:       false,
		MenuItems:    make([]*MenuItem, 0),
		font:         loader.LoadFont("./res/font.png", "./res/font.fnt", aspectRatio),
		SelectedItem: -1,
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
	for index, item := range m.MenuItems {
		if m.SelectedItem == index {
			guis = append(guis, item.guiTexture)
		} else {
			guis = append(guis, item.selectedTexture)
		}
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

func (m *Menu) ComputeSelectedItem(x, y float32) {
	m.SelectedItem = itemIndex(x, y, len(m.MenuItems))
}
