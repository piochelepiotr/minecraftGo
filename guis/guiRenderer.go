package guis

import (
    "github.com/piochelepiotr/minecraftGo/models"
    "github.com/piochelepiotr/minecraftGo/shaders"
    "github.com/piochelepiotr/minecraftGo/toolbox"
	"github.com/go-gl/gl/v4.1-core/gl"
)

type GuiRenderer struct {
    Quad models.RawModel
    Shader shaders.GuiShader
}


func (r *GuiRenderer) Render(guis []GuiTexture) {
    r.Shader.Program.Start()
    gl.BindVertexArray(r.Quad.VaoID)
    gl.EnableVertexAttribArray(0)
    gl.Enable(gl.BLEND)
    gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
    gl.Disable(gl.DEPTH_TEST)
    for _, gui := range guis {
        transformationMatrix := toolbox.CreateTransformationMatrix2D(gui.Position, gui.Scale)
        r.Shader.LoadTransformationMatrix(transformationMatrix)
        gl.ActiveTexture(gl.TEXTURE0)
        gl.BindTexture(gl.TEXTURE_2D, gui.Id)
        gl.DrawArrays(gl.TRIANGLE_STRIP, 0, r.Quad.VertexCount)
    }
    gl.Enable(gl.DEPTH_TEST)
    gl.Disable(gl.BLEND)
    gl.DisableVertexAttribArray(0)
    gl.BindVertexArray(0)
    r.Shader.Program.Stop()
}

func (r *GuiRenderer) CleanUp() {
    r.Shader.CleanUp()
}
