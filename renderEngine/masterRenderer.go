package renderEngine

import (
	"github.com/piochelepiotr/minecraftGo/shaders"
	"github.com/piochelepiotr/minecraftGo/entities"
	"github.com/piochelepiotr/minecraftGo/models"
)

type MasterRenderer struct {
    shader shaders.StaticShader
    renderer Renderer
    entities map[models.TexturedModel] []entities.Entity//keep an eye on the key
}

func CreateMasterRenderer() MasterRenderer {
    var r MasterRenderer
    r.shader = shaders.CreateStaticShader()
    r.renderer = CreateRenderer(r.shader)
    r.entities = make(map[models.TexturedModel] []entities.Entity)
    return r
}

func (r *MasterRenderer) Render(sun entities.Light, camera entities.Camera) {
    r.renderer.Prepare()
    r.shader.Program.Start()
    r.shader.LoadLight(sun)
    r.shader.LoadViewMatrix(camera)
    r.renderer.Render(r.entities)
    r.shader.Program.Stop()
    for model := range r.entities {
        delete(r.entities, model)
    }
}

func (r *MasterRenderer) ProcessEntity(entity entities.Entity) {
    r.entities[entity.TexturedModel] = append(r.entities[entity.TexturedModel], entity)
}

func (r *MasterRenderer) ProcessEntities(entities []entities.Entity) {
    for _, entity := range entities {
        r.ProcessEntity(entity)
    }
}

func (r *MasterRenderer) CleanUp() {
    r.shader.CleanUp()
}

