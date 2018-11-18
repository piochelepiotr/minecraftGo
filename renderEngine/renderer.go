package renderEngine

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/piochelepiotr/minecraftGo/toolbox"
	"github.com/piochelepiotr/minecraftGo/shaders"
	"github.com/piochelepiotr/minecraftGo/entities"
	"github.com/piochelepiotr/minecraftGo/models"
	"github.com/go-gl/mathgl/mgl32"
)

const fov = 70.0
const near_plane = 0.1
const far_plane = 1000

type Renderer struct {
    projectionMatrix mgl32.Mat4
    shader shaders.StaticShader
}

func CreateRenderer(shader shaders.StaticShader) Renderer {
    var r Renderer
    //gl.Enable(gl.CULL_FACE)
    //gl.CullFace(gl.BACK)
    r.projectionMatrix = mgl32.Perspective(mgl32.DegToRad(fov), float32(800.0)/600.0, near_plane, far_plane)
    r.shader = shader
    shader.Program.Start()
    shader.LoadProjectionMatrix(r.projectionMatrix)
    shader.Program.Stop()
    return r
}

func (r *Renderer) Prepare() {
    gl.Enable(gl.DEPTH_TEST)
    gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
    gl.ClearColor(0.5, 0.5, 0, 1)
}

func (r *Renderer) prepareTexturedModel(model models.TexturedModel) {
    gl.BindVertexArray(model.RawModel.VaoID)
    gl.EnableVertexAttribArray(0)
    gl.EnableVertexAttribArray(1)
    gl.EnableVertexAttribArray(2)
    r.shader.LoadShineVariables(model.ModelTexture.ShineDamper, model.ModelTexture.Reflectivity)
    gl.ActiveTexture(gl.TEXTURE0)
    gl.BindTexture(gl.TEXTURE_2D, model.ModelTexture.Id)
}

func (r *Renderer) unbindTexturedModel() {
    gl.DisableVertexAttribArray(0)
    gl.DisableVertexAttribArray(1)
    gl.DisableVertexAttribArray(2)
    gl.BindVertexArray(0)
}

func (r *Renderer) prepareEntity(entity entities.Entity) {
    transformationMatrix := toolbox.CreateTransformationMatrix(entity.Position, entity.Rotation, 1)
    r.shader.LoadTransformationMatrix(transformationMatrix)
}

func (r *Renderer) Render(allEntities map[models.TexturedModel] []entities.Entity) {
    for model := range allEntities {
        r.prepareTexturedModel(model)
        for _, entity := range allEntities[model] {
            r.prepareEntity(entity)
            gl.DrawElements(gl.TRIANGLES, model.RawModel.VertexCount, gl.UNSIGNED_INT, gl.PtrOffset(0))
        }
        r.unbindTexturedModel()
    }
}

