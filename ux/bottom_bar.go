package ux

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/piochelepiotr/minecraftGo/entities"
	"github.com/piochelepiotr/minecraftGo/game_engine/guis"
	"github.com/piochelepiotr/minecraftGo/game_engine/loader"
	"github.com/piochelepiotr/minecraftGo/game_engine/models"
	"github.com/piochelepiotr/minecraftGo/game_engine/render"
	pworld "github.com/piochelepiotr/minecraftGo/world"
	"github.com/piochelepiotr/minecraftGo/world/block"
)

const bottomBarHeight float32 = 0.1
const bottomBarItems = 9
const selectedSizeFactor float32 = 34.0/32.0

type BottomBar struct {
	background      guis.GuiTexture
	itemTextureID uint32
	selectItemTextureID uint32
	items []guis.GuiTexture
	selectedItem int
	aspectRatio float32
	objectsGuis []entities.Gui3dEntity
	content *pworld.Inventory
}

func NewBottomBar(aspectRatio float32, loader *loader.Loader, content *pworld.Inventory) *BottomBar {
	selectedItemTextureID := loader.LoadGuiTexture("textures/selected_item.png", mgl32.Vec2{}, mgl32.Vec2{}).Id
	itemTextureID := loader.LoadGuiTexture("textures/item.png", mgl32.Vec2{}, mgl32.Vec2{}).Id
	items := make([]guis.GuiTexture, bottomBarItems)
	b := &BottomBar{
		selectedItem: 0,
		background:      loader.LoadGuiTexture("textures/black.png", mgl32.Vec2{}, mgl32.Vec2{}),
		items: items,
		itemTextureID: itemTextureID,
		selectItemTextureID: selectedItemTextureID,
		aspectRatio: aspectRatio,
		content: content,
	}
	b.buildObjectsGuis(loader)
	b.selectItem(2)
	return b
}

func (b *BottomBar) GetSelectedBlock() block.Block {
	return b.content.BottomBar()[b.selectedItem].B
}

func (b *BottomBar) buildObjectsGuis(loader *loader.Loader) {
	modelTexture := loader.LoadModelTexture("textures/textures2.png")
	b.objectsGuis = make([]entities.Gui3dEntity, bottomBarItems)
	for i, o := range b.content.BottomBar() {
		if o.B != block.Air {
			model := models.TexturedModel{
				RawModel:  pworld.GetIconBlock(o.B, loader),
				TextureID: modelTexture,
			}
			b.objectsGuis[i] = entities.Gui3dEntity{
				Entity: entities.Entity{TexturedModel: model},
				Scale: bottomBarHeight,
			}
		}
	}
}

func (b *BottomBar) OffsetSelectedItem(offset int) {
	b.selectedItem += offset
	b.selectedItem %= bottomBarItems
	if b.selectedItem < 0 {
		b.selectedItem += bottomBarItems
	}
	b.selectItem(b.selectedItem)
}

func (b *BottomBar) selectItem(i int) {
	b.selectedItem = i
	b.buildItems()
}

func (b *BottomBar) Resize(aspectRatio float32) {
	b.aspectRatio = aspectRatio
	b.buildItems()
}

func (b *BottomBar) buildItems() {
	itemWidth := bottomBarHeight / b.aspectRatio
	bottomBarWidth := float32(bottomBarItems) * itemWidth
	posY := 1-bottomBarHeight
	b.background.Scale = mgl32.Vec2{bottomBarWidth, bottomBarHeight}
	b.background.Position = mgl32.Vec2{0, posY}
	for i := range b.items {
		b.items[i].Scale = mgl32.Vec2{itemWidth, bottomBarHeight}
		b.items[i].Position = mgl32.Vec2{float32(i-bottomBarItems/2)*2 * itemWidth, posY}
		b.objectsGuis[i].Translation = mgl32.Vec2{float32(i-bottomBarItems/2)*2 * itemWidth, posY}
		if i == b.selectedItem {
			b.items[i].Id = b.selectItemTextureID
			s := b.items[i].Scale
			b.items[i].Scale = mgl32.Vec2{s.X()*selectedSizeFactor, s.Y()*selectedSizeFactor}
		} else {
			b.items[i].Id = b.itemTextureID
		}
	}
}

// Render renders all objects on the screen
func (b *BottomBar) Render(renderer *render.MasterRenderer) {
	renderer.ProcessGui(b.background)
	for i, item := range b.items {
		if i != b.selectedItem {
			renderer.ProcessGui(item)
		}
	}
	renderer.ProcessGui(b.items[b.selectedItem])
	objects := b.content.BottomBar()
	for i, gui := range b.objectsGuis {
		if objects[i].B != block.Air {
			renderer.Process3DGui(gui)
		}
	}
}