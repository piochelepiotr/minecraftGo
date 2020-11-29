package world

import (
	"github.com/piochelepiotr/minecraftGo/world/block"
	"log"
)

const (
	MainItemsX = 9
	MainItemsY = 3
	BottomBar = 9
	Craft = 4
	CraftResult = 1
	maxItems = 64
	moving = 1
)
// mainItems + bottom bar + craft + craft result
const inventorySize = MainItemsX * MainItemsY + BottomBar + Craft + CraftResult + moving

type Item struct {
	B block.Block `json:"block"`
	N int `json:"quantity"`
}

type Inventory struct {
	Items []Item `json:"items"`
	mainItemsOffset int
	bottomBarOffset int
	craftOffset int
	craftResultOffset int
	movingOffset int
}

func (i *Inventory) BottomBar() []Item {
	return i.Items[i.bottomBarOffset: i.bottomBarOffset+BottomBar]
}

func (i *Inventory) MainItems() []Item {
	return i.Items[i.mainItemsOffset:i.mainItemsOffset + MainItemsX * MainItemsY]
}

func (i *Inventory) RemoveBottomBar(j int) {
	bar := i.BottomBar()
	bar[j].N -= 1
	if bar[j].N == 0 {
		bar[j].B = block.Air
	}
}

func (i *Inventory) GetMoving() (index int, b block.Block) {
	return i.movingOffset, i.Items[i.movingOffset].B
}

func (i *Inventory) Add(newB block.Block) {
	// 1. Tries to add it to the bottom bar on the same item
	bar := i.BottomBar()
	mainItems := i.MainItems()
	for i, b := range bar {
		if b.B == newB && b.N < maxItems {
			bar[i].N++
			return
		}
	}
	// 2. Free spot in the bottom bar
	for i, b := range bar {
		if b.B == block.Air {
			bar[i].B = newB
			bar[i].N = 1
			return
		}
	}
	// 3. Stack on same block in main items
	for i, b := range mainItems {
		if b.B == newB && b.N < maxItems {
			mainItems[i].N++
			return
		}
	}
	// 4. Find empty spot in main items
	for i, b := range mainItems {
		if b.B == block.Air {
			mainItems[i].B = newB
			mainItems[i].N = 1
			return
		}
	}
	log.Print("Full inventory")
}

func NewInventory() *Inventory {
	i := &Inventory{
		Items: make([]Item, inventorySize),
	}
	i.mainItemsOffset = 0
	i.bottomBarOffset = i.mainItemsOffset + MainItemsX*MainItemsY
	i.craftOffset = i.bottomBarOffset + BottomBar
	i.craftResultOffset = i.craftOffset + Craft
	i.movingOffset = i.craftResultOffset + 1
	for j := 0; j < inventorySize; j++ {
		i.Items[j] = Item{B: block.Air, N: 0}
	}
	// i.Items[len(i.Items)-1] = Item{B: block.Iron, N: 1}
	bar := i.BottomBar()
	bar[0] = Item{B: block.Cactus, N: 2}
	bar[1] = Item{B: block.Dirt, N: 14}
	bar[2] = Item{B: block.Sand, N: 53}
	bar[3] = Item{B: block.Stone, N: 64}
	bar[4] = Item{B: block.BirchLeaves, N: 1}
	bar[5] = Item{B: block.Birch, N: 1}
	bar[6] = Item{B: block.TallGrass, N: 1}
	i.Items[i.movingOffset] = Item{B: block.Iron, N:2}
	return i
}