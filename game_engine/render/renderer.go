package render

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/piochelepiotr/minecraftGo/entities"
	"github.com/piochelepiotr/minecraftGo/game_engine/models"
	"github.com/piochelepiotr/minecraftGo/game_engine/shaders"
	"github.com/piochelepiotr/minecraftGo/toolbox"
)

const nearPlane = 0.1
const farPlane = 1000

var skyColor = mgl32.Vec3{130.0/255.0, 166.0/255.0, 255.0/255.0}

var Fov = mgl32.DegToRad(70.0)

type Renderer struct {
	projectionMatrix mgl32.Mat4
	shader           shaders.StaticShader
}

func CreateRenderer() Renderer {
	shader := shaders.CreateStaticShader()
	EnableCulling()
	var r Renderer
	r.projectionMatrix = mgl32.Perspective(Fov, float32(800.0)/600.0, nearPlane, farPlane)
	r.shader = shader
	shader.Program.Start()
	shader.LoadProjectionMatrix(r.projectionMatrix)
	shader.Program.Stop()
	return r
}

func EnableCulling() {
	gl.Enable(gl.CULL_FACE)
	gl.CullFace(gl.BACK)
}

func DisableCulling() {
	gl.Disable(gl.CULL_FACE)
}

func (r *Renderer) Prepare() {
	gl.Enable(gl.DEPTH_TEST)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.ClearColor(skyColor[0], skyColor[1], skyColor[2], 1)
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
	r.shader.LoadShineVariables(model.ModelTexture.ShineDamper, model.ModelTexture.Reflectivity)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, model.ModelTexture.Id)
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

func (r *Renderer) Render(allEntities map[models.TexturedModel][]entities.Entity, light *entities.Light, camera *entities.Camera) {
	r.shader.Program.Start()
	defer r.shader.Program.Stop()
	if light != nil {
		r.shader.LoadLight(light)
	}
	if camera != nil {
		r.shader.LoadViewMatrix(camera)
	}
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
func (r *Renderer) CleanUp() {
		r.shader.CleanUp()
}
