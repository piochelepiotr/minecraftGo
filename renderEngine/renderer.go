package renderEngine

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/piochelepiotr/minecraftGo/models"
)

func Prepare() {
    gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
    gl.ClearColor(0.5, 0.5, 0, 1)
}

func Render(texturedModel models.TexturedModel) {
    model := texturedModel.RawModel
    gl.BindVertexArray(model.VaoID)
    gl.EnableVertexAttribArray(0)
    gl.EnableVertexAttribArray(1)
    gl.ActiveTexture(gl.TEXTURE0)
    gl.BindTexture(gl.TEXTURE_2D, texturedModel.ModelTexture.Id)
    gl.DrawElements(gl.TRIANGLES, model.VertexCount, gl.UNSIGNED_INT, gl.PtrOffset(0))
    gl.DisableVertexAttribArray(0)
    gl.DisableVertexAttribArray(1)
    gl.BindVertexArray(0)
}
