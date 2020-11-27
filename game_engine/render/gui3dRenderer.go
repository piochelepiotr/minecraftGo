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
const guiFarPlane = 100


var gui3dFov = mgl32.DegToRad(70.0)

type gui3dRenderer struct {
	projectionMatrix mgl32.Mat4
	shader           shaders.Gui3dShader
	camera entities.Camera
}

func createGui3dRenderer(aspectRatio float32) gui3dRenderer {
	shader := shaders.CreateGui3dShader()
	var r gui3dRenderer
	r.shader = shader
	r.camera = entities.Camera{Position: mgl32.Vec3{0, 0, 0}, Rotation: mgl32.Vec3{0, 0, 0}}
	r.shader.Program.Start()
	r.shader.LoadViewMatrix(&r.camera)
	r.shader.Program.Stop()
	r.resize(aspectRatio)
	return r
}

func (r *gui3dRenderer) resize(aspectRatio float32) {
	r.projectionMatrix = mgl32.Perspective(gui3dFov, aspectRatio, guiNearPlane, guiFarPlane)
	r.shader.Program.Start()
	r.shader.LoadProjectionMatrix(r.projectionMatrix)
	r.shader.LoadAspectRatio(aspectRatio)
	r.shader.Program.Stop()
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
	gl.BindTexture(gl.TEXTURE_2D, model.TextureID)
}

func (r *gui3dRenderer) unbindTexturedModel() {
	EnableCulling()
	gl.DisableVertexAttribArray(0)
	gl.DisableVertexAttribArray(1)
	gl.DisableVertexAttribArray(2)
	gl.DisableVertexAttribArray(3)
	gl.BindVertexArray(0)
}

func (r *gui3dRenderer) prepareEntity(entity entities.Gui3dEntity) {
	transformationMatrix := toolbox.CreateTransformationMatrix(mgl32.Vec3{}, mgl32.Vec3{0.5, 0.5, 0.25}, entity.Scale)
	r.shader.LoadTransformationMatrix(transformationMatrix)
	r.shader.LoadTranslation(entity.Translation)
}

func (r *gui3dRenderer) render(allEntities map[models.TexturedModel][]entities.Gui3dEntity) {
	// clear depth buffer before displaying 3d guis or they could be rendered behind the world
	gl.Clear(gl.DEPTH_BUFFER_BIT)
	r.shader.Program.Start()
	defer r.shader.Program.Stop()
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
