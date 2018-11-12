package renderEngine

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/piochelepiotr/minecraftGo/toolbox"
	"github.com/piochelepiotr/minecraftGo/shaders"
	"github.com/piochelepiotr/minecraftGo/entities"
	"github.com/go-gl/mathgl/mgl32"
)

const fov = 45.0
const near_plane = 0.1
const far_plane = 10.0
var projectionMatrix mgl32.Mat4

func init() {
    projectionMatrix = mgl32.Perspective(mgl32.DegToRad(fov), float32(800.0)/600.0, 0.1, 10.0)
}

func Prepare() {
    gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
    gl.ClearColor(0.5, 0.5, 0, 1)
}

func Render(player entities.Player, shader shaders.StaticShader) {
    texturedModel := player.TexturedModel
    model := texturedModel.RawModel
    gl.BindVertexArray(model.VaoID)
    gl.EnableVertexAttribArray(0)
    gl.EnableVertexAttribArray(1)
    transformationMatrix := toolbox.CreateTransformationMatrix(player.Entity.Position, player.Entity.Rotation, 1)
    shader.LoadTransformationMatrix(transformationMatrix)
    shader.LoadProjectionMatrix(projectionMatrix)
    gl.ActiveTexture(gl.TEXTURE0)
    gl.BindTexture(gl.TEXTURE_2D, texturedModel.ModelTexture.Id)
    gl.DrawElements(gl.TRIANGLES, model.VertexCount, gl.UNSIGNED_INT, gl.PtrOffset(0))
    gl.DisableVertexAttribArray(0)
    gl.DisableVertexAttribArray(1)
    gl.BindVertexArray(0)
}
