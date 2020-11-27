package render

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/piochelepiotr/minecraftGo/entities"
	"github.com/piochelepiotr/minecraftGo/game_engine/models"
	"github.com/piochelepiotr/minecraftGo/game_engine/shaders"
	"github.com/piochelepiotr/minecraftGo/toolbox"
	"sort"
)

const nearPlane = 0.1
const farPlane = 1000

var skyColor = mgl32.Vec3{130.0/255.0, 166.0/255.0, 255.0/255.0}

var Fov = mgl32.DegToRad(70.0)

type Renderer struct {
	projectionMatrix mgl32.Mat4
	shader           shaders.StaticShader
}

func CreateRenderer(aspectRatio float32) Renderer {
	shader := shaders.CreateStaticShader()
	EnableCulling()
	var r Renderer
	r.shader = shader
	r.resize(aspectRatio)
	return r
}

func (r *Renderer) resize(aspectRatio float32) {
	r.projectionMatrix = mgl32.Perspective(Fov, aspectRatio, nearPlane, farPlane)
	r.shader.Program.Start()
	r.shader.LoadProjectionMatrix(r.projectionMatrix)
	r.shader.Program.Stop()
}

func EnableCulling() {
	gl.Enable(gl.CULL_FACE)
	gl.CullFace(gl.BACK)
}

func DisableCulling() {
	gl.Disable(gl.CULL_FACE)
}

func (r *Renderer) prepareTexturedModel(model models.TexturedModel) {
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

func (r *Renderer) unbindTexturedModel() {
	EnableCulling()
	gl.DisableVertexAttribArray(0)
	gl.DisableVertexAttribArray(1)
	gl.DisableVertexAttribArray(2)
	gl.DisableVertexAttribArray(3)
	gl.BindVertexArray(0)
}

func (r *Renderer) prepareEntity(entity entities.Entity) {
	transformationMatrix := toolbox.CreateTransformationMatrix(entity.Position, entity.Rotation, 1)
	r.shader.LoadTransformationMatrix(transformationMatrix)
}

func (r *Renderer) Render(allEntities []entities.Entity, camera *entities.Camera) {
	// we sort entities to render always in the same order transparent blocks, to avoid flickering
	sort.Sort(entities.Entities(allEntities))
	r.shader.Program.Start()
	defer r.shader.Program.Stop()
	if camera != nil {
		r.shader.LoadViewMatrix(camera)
	}
	for _, entity := range allEntities {
		r.prepareTexturedModel(entity.TexturedModel)
		r.prepareEntity(entity)
		gl.DrawElements(gl.TRIANGLES, entity.TexturedModel.RawModel.VertexCount, gl.UNSIGNED_INT, gl.PtrOffset(0))
		r.unbindTexturedModel()
	}
}
// CleanUp frees memory for the shader
func (r *Renderer) CleanUp() {
		r.shader.CleanUp()
}
