package ux

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/piochelepiotr/minecraftGo/entities"
	"github.com/piochelepiotr/minecraftGo/game_engine/font"
	"github.com/piochelepiotr/minecraftGo/game_engine/guis"
	"github.com/piochelepiotr/minecraftGo/game_engine/loader"
	"github.com/piochelepiotr/minecraftGo/game_engine/models"
	"github.com/piochelepiotr/minecraftGo/game_engine/render"
	pworld "github.com/piochelepiotr/minecraftGo/world"
	"github.com/piochelepiotr/minecraftGo/world/block"
	"strconv"
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
	font         *font.FontType
	quantities   []font.GUIText
	loader *loader.Loader
	textureID uint32
}

func NewBottomBar(aspectRatio float32, loader *loader.Loader, content *pworld.Inventory) *BottomBar {
	selectedItemTextureID := loader.LoadGuiTexture("textures/selected_item.png", mgl32.Vec2{}, mgl32.Vec2{}).Id
	itemTextureID := loader.LoadGuiTexture("textures/item.png", mgl32.Vec2{}, mgl32.Vec2{}).Id
	items := make([]guis.GuiTexture, bottomBarItems)
	b := &BottomBar{
		font:         loader.LoadFont("./res/font.png", "./res/font.fnt", aspectRatio),
		selectedItem: 0,
		background:      loader.LoadGuiTexture("textures/black.png", mgl32.Vec2{}, mgl32.Vec2{}),
		items: items,
		itemTextureID: itemTextureID,
		selectItemTextureID: selectedItemTextureID,
		aspectRatio: aspectRatio,
		content: content,
		quantities: make([]font.GUIText, len(content.BottomBar())),
		objectsGuis: make([]entities.Gui3dEntity, len(content.BottomBar())),
		textureID: loader.LoadModelTexture("textures/textures2.png"),
		loader: loader,
	}
	b.ReBuild()
	return b
}

func (b *BottomBar) GetSelectedBlock() block.Block {
	return b.content.BottomBar()[b.selectedItem].B
}

func (b *BottomBar) ReBuild() {
	for i, o := range b.content.BottomBar() {
		if o.B != block.Air {
			model := models.TexturedModel{
				RawModel:  pworld.GetIconBlock(o.B, b.loader),
				TextureID: b.textureID,
			}
			b.objectsGuis[i].Entity.TexturedModel = model
			if o.N > 1 {
				b.quantities[i] = b.loader.LoadText(font.CreateGUIText(strconv.Itoa(o.N), 1.4, b.font, mgl32.Vec2{0, 0}, 1, 1, false, mgl32.Vec3{1, 1, 1}))
			}
		}
	}
	b.buildItems()
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
	b.buildItems()
}

func (b *BottomBar) Resize(aspectRatio float32) {
	b.aspectRatio = aspectRatio
	b.buildItems()
}

func (b *BottomBar) updatePos(item int, width, height, posX, posY, itemWidth float32) {
	b.items[item].Scale = mgl32.Vec2{width, height}
	b.items[item].Position = mgl32.Vec2{posX, posY}
	b.objectsGuis[item].Translation = mgl32.Vec2{posX, posY}
	b.objectsGuis[item].Scale = bottomBarHeight
	text := b.quantities[item]
	b.quantities[item].Position = mgl32.Vec2{posX+pos(itemWidth/2-text.Width)+pos(0.007), posY+pos(bottomBarHeight/2-text.GetLineHeight()+pos(0.005))}
}

func (b *BottomBar) buildItems() {
	itemWidth := bottomBarHeight / b.aspectRatio
	bottomBarWidth := float32(bottomBarItems) * itemWidth
	posY := pos(0.5-bottomBarHeight/2)
	startX := pos(-bottomBarWidth/2+itemWidth/2)
	b.background.Scale = mgl32.Vec2{bottomBarWidth, bottomBarHeight}
	b.background.Position = mgl32.Vec2{0, posY}
	for i := range b.items {
		b.updatePos(i, itemWidth, bottomBarHeight, startX + pos(itemWidth * float32(i)), posY, itemWidth)
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
	for i, quantity := range b.quantities {
		if objects[i].N > 1 {
			renderer.ProcessText(quantity)
		}
	}
}