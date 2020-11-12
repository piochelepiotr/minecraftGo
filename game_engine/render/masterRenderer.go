package render

import (
	"github.com/piochelepiotr/minecraftGo/entities"
	"github.com/piochelepiotr/minecraftGo/game_engine/font"
	pguis "github.com/piochelepiotr/minecraftGo/game_engine/guis"
	"github.com/piochelepiotr/minecraftGo/game_engine/loader"
	"github.com/piochelepiotr/minecraftGo/game_engine/models"
	"github.com/piochelepiotr/minecraftGo/game_engine/shaders"
)

// MasterRenderer is the main renderer that will render
// everything on the screen
type MasterRenderer struct {
	light        *entities.Light
	camera       *entities.Camera
	shader       shaders.StaticShader
	renderer     Renderer
	fontRenderer FontRenderer
	guiRenderer  pguis.GuiRenderer
	entities     map[models.TexturedModel][]entities.Entity //keep an eye on the key
	guis         []pguis.GuiTexture
	texts        []font.GUIText
}

// CreateMasterRenderer creates a MasterRenderer class
func CreateMasterRenderer() *MasterRenderer {
	var r MasterRenderer
	r.fontRenderer = CreateFontRenderer()
	r.guiRenderer = loader.CreateGuiRenderer()
	r.shader = shaders.CreateStaticShader()
	r.renderer = CreateRenderer(r.shader)
	r.entities = make(map[models.TexturedModel][]entities.Entity)
	r.guis = make([]pguis.GuiTexture, 0)
	r.texts = make([]font.GUIText, 0)
	return &r
}

// Render renders everything on the screen
func (r *MasterRenderer) Render() {
	r.renderer.Prepare()
	r.shader.Program.Start()
	if r.light != nil {
		r.shader.LoadLight(r.light)
	}
	if r.camera != nil {
		r.shader.LoadViewMatrix(r.camera)
	}
	r.renderer.Render(r.entities)
	r.shader.Program.Stop()
	r.guiRenderer.Render(r.guis)
	r.fontRenderer.Render(r.texts)
	for model := range r.entities {
		delete(r.entities, model)
	}
	r.guis = make([]pguis.GuiTexture, 0)
	r.texts = make([]font.GUIText, 0)
	r.light = nil
	r.camera = nil
}

func (r *MasterRenderer) SetLight(light *entities.Light) {
	r.light = light
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

// CleanUp frees memory for the shader and the renderers
func (r *MasterRenderer) CleanUp() {
	r.shader.CleanUp()
	r.fontRenderer.CleanUp()
	r.guiRenderer.CleanUp()
}
