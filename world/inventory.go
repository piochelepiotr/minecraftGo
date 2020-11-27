package world

import "github.com/piochelepiotr/minecraftGo/world/block"

const (
	MainItemsX = 11
	MainItemsY = 6
	BottomBar = 9
	Craft = 4
	CraftResult = 1
)
// mainItems + bottom bar + craft + craft result
const inventorySize = MainItemsX * MainItemsY + BottomBar + Craft + CraftResult

type Item struct {
	B block.Block `json:"block"`
	N int `json:"quantity"`
}

type Inventory struct {
	Items []Item `json:"items"`
}

func (i *Inventory) BottomBar() []Item {
	return i.Items[MainItemsX * MainItemsY: MainItemsX*MainItemsY+BottomBar]
}

func NewInventory() *Inventory {
	i := &Inventory{
		Items: make([]Item, inventorySize),
	}
	for j := 0; j < inventorySize; j++ {
		i.Items[j] = Item{B: block.Grass, N: 1}
	}
	i.Items[len(i.Items)-1] = Item{B: block.Iron, N: 1}
	bar := i.BottomBar()
	bar[0] = Item{B: block.Cactus, N: 2}
	bar[1] = Item{B: block.Dirt, N: 1}
	bar[2] = Item{B: block.Sand, N: 1}
	bar[3] = Item{B: block.Stone, N: 1}
	bar[4] = Item{B: block.BirchLeaves, N: 1}
	bar[5] = Item{B: block.Birch, N: 1}
	bar[6] = Item{B: block.TallGrass, N: 1}
	return i
}