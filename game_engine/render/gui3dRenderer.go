package render

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/piochelepiotr/minecraftGo/entities"
	"github.com/piochelepiotr/minecraftGo/game_engine/models"
	"github.com/piochelepiotr/minecraftGo/game_engine/shaders"
	"github.com/piochelepiotr/minecraftGo/toolbox"
)

const guiNearPlane = 0.1
const guiFarPlane = 1000


var gui3dFov = mgl32.DegToRad(70.0)

type gui3dRenderer struct {
	projectionMatrix mgl32.Mat4
	shader           shaders.Gui3dShader
}

func createGui3dRenderer() gui3dRenderer {
	shader := shaders.CreateGui3dShader()
	var r gui3dRenderer
	r.projectionMatrix = mgl32.Perspective(Fov, float32(800.0)/600.0, guiNearPlane, guiFarPlane)
	r.shader = shader
	shader.Program.Start()
	shader.LoadProjectionMatrix(r.projectionMatrix)
	shader.Program.Stop()
	return r
}

func (r *gui3dRenderer) prepareTexturedModel(model models.TexturedModel) {
	if model.Transparent {
		DisableCulling()
	}
	gl.BindVertexArray(model.RawModel.VaoID)
	gl.EnableVertexAttribArray(0)
	gl.EnableVertexAttribArray(1)
	gl.EnableVertexAttribArray(2)
	gl.EnableVertexAttribArray(3)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, model.ModelTexture.Id)
}

func (r *gui3dRenderer) unbindTexturedModel() {
	EnableCulling()
	gl.DisableVertexAttribArray(0)
	gl.DisableVertexAttribArray(1)
	gl.DisableVertexAttribArray(2)
	gl.DisableVertexAttribArray(3)
	gl.BindVertexArray(0)
}

func (r *gui3dRenderer) prepareEntity(entity entities.Entity) {
	transformationMatrix := toolbox.CreateTransformationMatrix(entity.Position, entity.Rotation, 1)
	r.shader.LoadTransformationMatrix(transformationMatrix)
}

func (r *gui3dRenderer) render(allEntities map[models.TexturedModel][]entities.Entity) {
	// clear depth buffer before displaying 3d guis or they could be rendered behind the world
	gl.Clear(gl.DEPTH_BUFFER_BIT)
	r.shader.Program.Start()
	defer r.shader.Program.Stop()
	camera := entities.Camera{Position: mgl32.Vec3{0, 0, 0}, Rotation: mgl32.Vec3{0, 0, 0}}
	r.shader.LoadViewMatrix(&camera)
	for model := range allEntities {
		r.prepareTexturedModel(model)
		for _, entity := range allEntities[model] {
			r.prepareEntity(entity)
			gl.DrawElements(gl.TRIANGLES, model.RawModel.VertexCount, gl.UNSIGNED_INT, gl.PtrOffset(0))
		}
		r.unbindTexturedModel()
	}
}
// CleanUp frees memory for the shader
func (r *gui3dRenderer) cleanUp() {
		r.shader.CleanUp()
}
