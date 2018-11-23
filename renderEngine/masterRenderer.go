package renderEngine

import (
	"github.com/piochelepiotr/minecraftGo/entities"
	"github.com/piochelepiotr/minecraftGo/fontMeshCreator"
	"github.com/piochelepiotr/minecraftGo/fontRendering"
	pguis "github.com/piochelepiotr/minecraftGo/guis"
	"github.com/piochelepiotr/minecraftGo/loader"
	"github.com/piochelepiotr/minecraftGo/menu"
	"github.com/piochelepiotr/minecraftGo/models"
	"github.com/piochelepiotr/minecraftGo/shaders"
)

type MasterRenderer struct {
	shader       shaders.StaticShader
	renderer     Renderer
	fontRenderer fontRendering.FontRenderer
	guiRenderer  pguis.GuiRenderer
	entities     map[models.TexturedModel][]entities.Entity //keep an eye on the key
	guis         []pguis.GuiTexture
	texts        []fontMeshCreator.GUIText
}

func CreateMasterRenderer() MasterRenderer {
	var r MasterRenderer
	r.fontRenderer = fontRendering.CreateFontRenderer()
	r.guiRenderer = loader.CreateGuiRenderer()
	r.shader = shaders.CreateStaticShader()
	r.renderer = CreateRenderer(r.shader)
	r.entities = make(map[models.TexturedModel][]entities.Entity)
	r.guis = make([]pguis.GuiTexture, 0)
	r.texts = make([]fontMeshCreator.GUIText, 0)
	return r
}

func (r *MasterRenderer) Render(sun entities.Light, camera entities.Camera) {
	r.renderer.Prepare()
	r.shader.Program.Start()
	r.shader.LoadLight(sun)
	r.shader.LoadViewMatrix(camera)
	r.renderer.Render(r.entities)
	r.shader.Program.Stop()
	r.guiRenderer.Render(r.guis)
	r.fontRenderer.Render(r.texts)
	for model := range r.entities {
		delete(r.entities, model)
	}
	r.guis = make([]pguis.GuiTexture, 0)
	r.texts = make([]fontMeshCreator.GUIText, 0)
}

func (r *MasterRenderer) ProcessEntity(entity entities.Entity) {
	r.entities[entity.TexturedModel] = append(r.entities[entity.TexturedModel], entity)
}

func (r *MasterRenderer) ProcessEntities(entities []entities.Entity) {
	for _, entity := range entities {
		r.ProcessEntity(entity)
	}
}

func (r *MasterRenderer) ProcessGui(gui pguis.GuiTexture) {
	r.guis = append(r.guis, gui)
}

func (r *MasterRenderer) ProcessGuis(guis []pguis.GuiTexture) {
	r.guis = append(r.guis, guis...)
}

func (r *MasterRenderer) ProcessMenu(menu menu.Menu) {
	if menu.Opened {
		r.ProcessGuis(menu.GetMenuItems())
		r.ProcessTexts(menu.GetMenuTexts())
	}
}

func (r *MasterRenderer) ProcessText(text fontMeshCreator.GUIText) {
	r.texts = append(r.texts, text)
}

func (r *MasterRenderer) ProcessTexts(texts []fontMeshCreator.GUIText) {
	r.texts = append(r.texts, texts...)
}

func (r *MasterRenderer) CleanUp() {
	r.shader.CleanUp()
	r.fontRenderer.CleanUp()
	r.guiRenderer.CleanUp()
}
