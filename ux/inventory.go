package ux

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/piochelepiotr/minecraftGo/entities"
	"github.com/piochelepiotr/minecraftGo/game_engine/font"
	"github.com/piochelepiotr/minecraftGo/game_engine/guis"
	"github.com/piochelepiotr/minecraftGo/game_engine/loader"
	"github.com/piochelepiotr/minecraftGo/game_engine/models"
	"github.com/piochelepiotr/minecraftGo/game_engine/render"
	"github.com/piochelepiotr/minecraftGo/world"
	"github.com/piochelepiotr/minecraftGo/world/block"
	"strconv"
)

const (
	itemHeight = 0.08
	inventoryBorderHeight = 0.02
)

type Inventory struct {
	background      guis.GuiTexture
	aspectRatio float32
	content *world.Inventory
	itemsBackgrounds []guis.GuiTexture
	itemBackgroundTextureID uint32
	items3D []entities.Gui3dEntity
	items3DTextureID uint32
	font         *font.FontType
	quantities   []font.GUIText
	loader *loader.Loader
	selectedItem int
}

func NewInventory(aspectRatio float32, loader *loader.Loader, content *world.Inventory) *Inventory {
	itemBackgroundTextureID := loader.LoadGuiTexture("textures/inventory_item.png", mgl32.Vec2{}, mgl32.Vec2{}).Id

	i := &Inventory{
		font:         loader.LoadFont("./res/font.png", "./res/font.fnt", aspectRatio),
		background:              loader.LoadGuiTexture("textures/white.png", mgl32.Vec2{}, mgl32.Vec2{}),
		aspectRatio:             aspectRatio,
		content:                 content,
		itemsBackgrounds:        make([]guis.GuiTexture, len(content.Items)),
		itemBackgroundTextureID: itemBackgroundTextureID,
		items3DTextureID: loader.LoadModelTexture("textures/textures2.png"),
		items3D : make([]entities.Gui3dEntity, len(content.Items)),
		quantities: make([]font.GUIText, len(content.Items)),
		loader: loader,
		selectedItem: -1,
	}
	i.ReBuild()
	return i
}

func (i *Inventory) ReBuild() {
	for j, o := range i.content.Items {
		if o.B != block.Air {
			model := models.TexturedModel{
				RawModel:  world.GetIconBlock(o.B, i.loader),
				TextureID: i.items3DTextureID,
			}
			i.items3D[j] = entities.Gui3dEntity{
				Entity: entities.Entity{TexturedModel: model},
				Scale: itemHeight+0.01,
			}
			if o.N > 1 {
				i.quantities[j] = i.loader.LoadText(font.CreateGUIText(strconv.Itoa(o.N), 1.3, i.font, mgl32.Vec2{0, 0}, 1, 1, false, mgl32.Vec3{1, 1, 1}))
			}
		}
	}
	i.updatePositions()
}

func (i *Inventory) GetSelected() int {
	return i.selectedItem
}

func (i *Inventory) MoveMouse(x, y float32) {
	itemWidth := itemHeight / i.aspectRatio
	// moving object
	j, _ := i.content.GetMoving()
	i.updatePos(j, itemWidth, itemHeight, pos(x), pos(y))
	i.selectedItem = i.getSelectedItem(x, y)
}

// pos goes from -1 to 1, so we need to multiply everything by 2
func pos(x float32) float32 {
	return x*2
}

func (i *Inventory) updatePos(item int, width, height, posX, posY float32) {
	i.itemsBackgrounds[item].Scale = mgl32.Vec2{width, height}
	i.itemsBackgrounds[item].Position = mgl32.Vec2{posX, posY}
	i.items3D[item].Translation = mgl32.Vec2{posX, posY}
	text := i.quantities[item]
	i.quantities[item].Position = mgl32.Vec2{posX+pos(width/2-text.Width)+pos(0.007), posY+pos(itemHeight/2-text.GetLineHeight()+pos(0.005))}
}

func (i *Inventory) updatePositions() {
	itemWidth := itemHeight / i.aspectRatio
	borderWidth := inventoryBorderHeight / i.aspectRatio
	inventoryWidth := float32(world.MainItemsX) * itemWidth + borderWidth * 2
	inventoryHeight := float32(world.MainItemsY+3) * itemHeight + inventoryBorderHeight * 4
	i.background.Scale = mgl32.Vec2{inventoryWidth, inventoryHeight}

	startY := pos(-inventoryHeight/2)
	startX := pos(-inventoryWidth/2)

	// craft
	craftYOffset := startY + pos(inventoryBorderHeight + itemHeight/2)
	craftXOffset := startX + pos(borderWidth + itemWidth/2)
	for x := 0; x < 2; x++ {
		for y := 0; y < 2; y++ {
			j := world.MainItemsX*world.MainItemsY+bottomBarItems+y*2+x
			i.updatePos(j, itemWidth, itemHeight, pos(float32(x) * itemWidth)+ craftXOffset, pos(float32(y) * itemHeight) + craftYOffset)
			i.itemsBackgrounds[j].Id = i.itemBackgroundTextureID
		}
	}
	craftResultX := craftXOffset + pos(2 * itemWidth + borderWidth)
	craftResultY := craftYOffset + pos(itemHeight/2)
	j := world.MainItemsX*world.MainItemsY+bottomBarItems+world.Craft
	i.updatePos(j, itemWidth, itemHeight, craftResultX, craftResultY)
	i.itemsBackgrounds[j].Id = i.itemBackgroundTextureID

	// main items
	mainItemsYOffset := craftYOffset + pos(inventoryBorderHeight + itemHeight*2)
	mainItemsXOffset := startX + pos(borderWidth + itemWidth/2)
	for x := 0; x < world.MainItemsX; x++ {
		for y := 0; y < world.MainItemsY; y++ {
			j := y * world.MainItemsX + x
			i.updatePos(j, itemWidth, itemHeight, pos(float32(x) * itemWidth)+mainItemsXOffset, pos(float32(y) * itemHeight) + mainItemsYOffset)
			i.itemsBackgrounds[j].Id = i.itemBackgroundTextureID
		}
	}

	// bottom bar
	bottomBarY := mainItemsYOffset + pos(float32(world.MainItemsY) * itemHeight + inventoryBorderHeight)
	bottomBarXOffset := mainItemsXOffset + pos(float32(world.MainItemsX-world.BottomBar)*itemWidth/2)
	for x := 0; x < bottomBarItems; x++ {
		j := world.MainItemsX*world.MainItemsY+x
		i.updatePos(j, itemWidth, itemHeight, pos(float32(x) * itemWidth) + bottomBarXOffset, bottomBarY)
		i.itemsBackgrounds[j].Id = i.itemBackgroundTextureID
	}

	// moving object
	j, _ = i.content.GetMoving()
	p := i.itemsBackgrounds[j].Position
	i.updatePos(j, itemWidth, itemHeight, p.X(), p.Y())
}

func (i *Inventory) Resize(aspectRatio float32) {
	i.aspectRatio = aspectRatio
	i.updatePositions()
}

func abs(x float32) float32 {
	if x > 0 {
		return x
	}
	return -x
}

func (i *Inventory) getSelectedItem(x, y float32) int {
	x = pos(x)
	y = pos(y)
	itemWidth := itemHeight / i.aspectRatio
	w := pos(itemWidth)/2
	h := pos(itemHeight)/2
	movingIndex, _ := i.content.GetMoving()
	for j, item := range i.itemsBackgrounds {
		if j == movingIndex {
			continue
		}
		if abs(item.Position.X() - x) < w && abs(item.Position.Y() - y) < h {
			return j
		}
	}
	return -1
}

// Render renders all objects on the screen
func (i *Inventory) Render(renderer *render.MasterRenderer) {
	renderer.ProcessGui(i.background)
	movingIndex, _ := i.content.GetMoving()
	for j, item := range i.itemsBackgrounds {
		if j == movingIndex {
			continue
		}
		renderer.ProcessGui(item)
	}
	for j, item3D := range i.items3D {
		if i.content.Items[j].B != block.Air {
			renderer.Process3DGui(item3D)
		}
	}
	for j, quantity := range i.quantities {
		if i.content.Items[j].N > 1 {
			renderer.ProcessText(quantity)
		}
	}
}
