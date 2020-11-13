package render

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/piochelepiotr/minecraftGo/entities"
	"github.com/piochelepiotr/minecraftGo/game_engine/font"
	pguis "github.com/piochelepiotr/minecraftGo/game_engine/guis"
	"github.com/piochelepiotr/minecraftGo/game_engine/loader"
	"github.com/piochelepiotr/minecraftGo/game_engine/models"
)

// MasterRenderer is the main renderer that will render
// everything on the screen
type MasterRenderer struct {
	camera       *entities.Camera
	renderer     Renderer
	fontRenderer FontRenderer
	guiRenderer  pguis.GuiRenderer
	// gui3DRenderer is used to render cubes in the inventory
	gui3DRenderer gui3dRenderer
	entities     map[models.TexturedModel][]entities.Entity
	guis3D     map[models.TexturedModel][]entities.Entity
	guis         []pguis.GuiTexture
	texts        []font.GUIText
}

// CreateMasterRenderer creates a MasterRenderer class
func CreateMasterRenderer(aspectRatio float32) *MasterRenderer {
	var r MasterRenderer
	r.fontRenderer = CreateFontRenderer()
	r.guiRenderer = loader.CreateGuiRenderer()
	r.renderer = CreateRenderer(aspectRatio)
	r.gui3DRenderer = createGui3dRenderer()
	r.entities = make(map[models.TexturedModel][]entities.Entity)
	r.guis3D = make(map[models.TexturedModel][]entities.Entity)
	r.guis = make([]pguis.GuiTexture, 0)
	r.texts = make([]font.GUIText, 0)
	return &r
}

func (r *MasterRenderer) Resize(aspectRatio float32) {
	r.renderer.resize(aspectRatio)
}

func (r *MasterRenderer) Prepare() {
	gl.Enable(gl.DEPTH_TEST)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.ClearColor(skyColor[0], skyColor[1], skyColor[2], 1)
}

// Render renders everything on the screen
func (r *MasterRenderer) Render() {
	r.Prepare()
	r.renderer.Render(r.entities, r.camera)
	r.guiRenderer.Render(r.guis)
	r.fontRenderer.Render(r.texts)
	// r.gui3DRenderer.render(r.guis3D)
	r.entities = make(map[models.TexturedModel][]entities.Entity)
	r.guis3D = make(map[models.TexturedModel][]entities.Entity)
	r.guis = make([]pguis.GuiTexture, 0)
	r.texts = make([]font.GUIText, 0)
}

func (r *MasterRenderer) SetCamera(camera *entities.Camera) {
	r.camera = camera
}

// ProcessEntity adds entity to the list of entities to render
func (r *MasterRenderer) ProcessEntity(entity entities.Entity) {
	r.entities[entity.TexturedModel] = append(r.entities[entity.TexturedModel], entity)
}

// ProcessEntities adds entities to the list of entities to render
func (r *MasterRenderer) ProcessEntities(entities []entities.Entity) {
	for _, entity := range entities {
		r.ProcessEntity(entity)
	}
}

// ProcessGui adds gui to the list of guis to render
func (r *MasterRenderer) ProcessGui(gui pguis.GuiTexture) {
	r.guis = append(r.guis, gui)
}

// ProcessGuis adds guis to the list of guis to render
func (r *MasterRenderer) ProcessGuis(guis []pguis.GuiTexture) {
	r.guis = append(r.guis, guis...)
}

// ProcessText adds text to the list of texts to render
func (r *MasterRenderer) ProcessText(text font.GUIText) {
	r.texts = append(r.texts, text)
}

// ProcessTexts adds texts to the list of texts to render
func (r *MasterRenderer) ProcessTexts(texts []font.GUIText) {
	r.texts = append(r.texts, texts...)
}

func (r *MasterRenderer) Process3DGui(entity entities.Entity) {
	r.guis3D[entity.TexturedModel] = append(r.entities[entity.TexturedModel], entity)
}

// CleanUp frees memory for the shader and the renderers
func (r *MasterRenderer) CleanUp() {
	r.renderer.CleanUp()
	r.fontRenderer.CleanUp()
	r.guiRenderer.CleanUp()
	r.gui3DRenderer.cleanUp()
}
