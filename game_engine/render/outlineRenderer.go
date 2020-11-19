package render

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/piochelepiotr/minecraftGo/entities"
	"github.com/piochelepiotr/minecraftGo/game_engine/models"
	"github.com/piochelepiotr/minecraftGo/game_engine/shaders"
	"github.com/piochelepiotr/minecraftGo/toolbox"
)

type OutlineRenderer struct {
	projectionMatrix mgl32.Mat4
	shader           shaders.OutlineShader
}

func NewOutlineRenderer(aspectRatio float32) *OutlineRenderer {
	gl.Enable(gl.LINE_SMOOTH)
	r := &OutlineRenderer{
		shader: shaders.CreateOutlineShader(),
	}
	r.resize(aspectRatio)
	return r
}

func (r *OutlineRenderer) resize(aspectRatio float32) {
	r.projectionMatrix = mgl32.Perspective(Fov, aspectRatio, nearPlane, farPlane)
	r.shader.Program.Start()
	r.shader.LoadProjectionMatrix(r.projectionMatrix)
	r.shader.Program.Stop()
}

func (r *OutlineRenderer) prepareModel(model models.OutlineModel) {
	gl.BindVertexArray(model.RawModel.VaoID)
	gl.EnableVertexAttribArray(0)

	transformationMatrix := toolbox.CreateTransformationMatrix(model.Position, mgl32.Vec3{}, 1)
	r.shader.LoadTransformationMatrix(transformationMatrix)
}

func (r OutlineRenderer) unbindModel() {
	gl.DisableVertexAttribArray(0)
	gl.BindVertexArray(0)
}

func (r *OutlineRenderer) render(models []models.OutlineModel, camera *entities.Camera) {
	if len(models) == 0 {
		return
	}
	r.shader.Program.Start()
	defer r.shader.Program.Stop()
	r.shader.LoadViewMatrix(camera)
	for _, m := range models {
		r.prepareModel(m)
		gl.DrawElements(gl.LINES, m.VertexCount, gl.UNSIGNED_INT, gl.PtrOffset(0))
		r.unbindModel()
	}
}

// CleanUp frees memory for the shader
func (r *OutlineRenderer) cleanUp() {
	r.shader.CleanUp()
}
