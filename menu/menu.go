package menu

import (
    "github.com/piochelepiotr/minecraftGo/guis"
)

type Menu struct {
    Opened bool
    MenuItems []*MenuItem
}

func CreateMenu() Menu {
    return Menu{
        Opened: false,
        MenuItems: make([]*MenuItem, 0),
    }
}

func (m *Menu) AddItem(text string) {
    m.MenuItems = append(m.MenuItems, CreateMenuItem(text, len(m.MenuItems)))
    for _, item := range m.MenuItems {
        item.computeYPos(len(m.MenuItems))
    }
}

func (m *Menu) GetMenuItems() []guis.GuiTexture {
    guis := make([]guis.GuiTexture,0)
    for _, item := range m.MenuItems {
        guis = append(guis, item.guiTexture)
    }
    return guis
}
