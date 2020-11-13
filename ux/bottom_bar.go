package ux

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/piochelepiotr/minecraftGo/entities"
	"github.com/piochelepiotr/minecraftGo/game_engine/guis"
	"github.com/piochelepiotr/minecraftGo/game_engine/loader"
	"github.com/piochelepiotr/minecraftGo/game_engine/models"
	"github.com/piochelepiotr/minecraftGo/game_engine/render"
	pworld "github.com/piochelepiotr/minecraftGo/world"
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
	objects []pworld.Block
	objectsGuis []entities.Gui3dEntity
}

func NewBottomBar(aspectRatio float32) *BottomBar {
	selectedItemTextureID := loader.LoadGuiTexture("textures/selected_item.png", mgl32.Vec2{}, mgl32.Vec2{}).Id
	itemTextureID := loader.LoadGuiTexture("textures/item.png", mgl32.Vec2{}, mgl32.Vec2{}).Id
	items := make([]guis.GuiTexture, 0, bottomBarItems)
	objects := make([]pworld.Block, 0, bottomBarItems)
	for i := 0; i < bottomBarItems; i++ {
		objects = append(objects, pworld.Grass)
	}
	objects[0] = pworld.Cactus
	objects[1] = pworld.Dirt
	objects[2] = pworld.Sand
	objects[3] = pworld.Stone
	objects[4] = pworld.Tree
	for i := 0; i < bottomBarItems; i++ {
		items = append(items, guis.GuiTexture{})
	}
	b := &BottomBar{
		selectedItem: 0,
		background:      loader.LoadGuiTexture("textures/black.png", mgl32.Vec2{}, mgl32.Vec2{}),
		items: items,
		itemTextureID: itemTextureID,
		selectItemTextureID: selectedItemTextureID,
		aspectRatio: aspectRatio,
		objects: objects,
	}
	b.buildObjectsGuis()
	b.selectItem(2)
	return b
}

func (b *BottomBar) GetSelectedBlock() pworld.Block {
	return b.objects[b.selectedItem]
}

func (b *BottomBar) buildObjectsGuis() {
	modelTexture := loader.LoadModelTexture("textures/textures2.png", 16)
	b.objectsGuis = make([]entities.Gui3dEntity, bottomBarItems)
	for i, o := range b.objects {
		if o != pworld.Air {
			model := models.TexturedModel{
				RawModel:     getBlockModel(o),
				ModelTexture: modelTexture,
			}
			b.objectsGuis[i] = entities.Gui3dEntity{
				Entity: entities.Entity{TexturedModel: model},
			}
		}
	}
}

func getBlockModel(b pworld.Block) models.RawModel {
	block := pworld.NewChunk(pworld.NewOneBlockChunk(b))
	block.ChangeOrigin()
	block.Load()
	return block.Model
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
	for i, gui := range b.objectsGuis {
		if b.objects[i] != pworld.Air {
			renderer.Process3DGui(gui)
		}
	}
}

