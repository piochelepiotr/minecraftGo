package menu

import (
	"github.com/piochelepiotr/minecraftGo/font"
	"github.com/piochelepiotr/minecraftGo/guis"
	"github.com/piochelepiotr/minecraftGo/loader"
)

// Menu is the menu of the game
type Menu struct {
	Opened       bool
	Items        []*Item
	font         *font.FontType
	SelectedItem int
}

// CreateMenu creates the menu of the game
func CreateMenu(aspectRatio float32) Menu {
	return Menu{
		Opened:       false,
		Items:        make([]*Item, 0),
		font:         loader.LoadFont("./res/font.png", "./res/font.fnt", aspectRatio),
		SelectedItem: -1,
	}
}

// AddItem adds item to the game menu
func (m *Menu) AddItem(text string) {
	m.Items = append(m.Items, CreateItem(text, len(m.Items), m.font))
	for _, item := range m.Items {
		item.computeYPos(len(m.Items))
	}
}

// GetItems gets all guis of the menu
func (m *Menu) GetItems() []guis.GuiTexture {
	guis := make([]guis.GuiTexture, 0)
	for index, item := range m.Items {
		if m.SelectedItem == index {
			guis = append(guis, item.guiTexture)
		} else {
			guis = append(guis, item.selectedTexture)
		}
	}
	return guis
}

// GetMenuTexts gets all texts of the menu
func (m *Menu) GetMenuTexts() []font.GUIText {
	texts := make([]font.GUIText, 0)
	for _, item := range m.Items {
		texts = append(texts, item.text)
	}
	return texts
}

// ComputeSelectedItem returns the index of the item under the cursor
func (m *Menu) ComputeSelectedItem(x, y float32) {
	m.SelectedItem = itemIndex(x, y, len(m.Items))
}
