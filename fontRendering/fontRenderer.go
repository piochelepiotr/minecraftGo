package fontRendering

import (
    "github.com/piochelepiotr/minecraftGo/shaders"
	"github.com/go-gl/gl/v4.1-core/gl"
)


type FontRenderer struct {
	shader shaders.FontShader
    texts map[fontMeshCreator.FontType][]fontMeshCreator.GUIText
}

func CreateFontRenderer() FontRenderer {
    return FontRenderer{
	    shader: shaders.CreateFontShader()
        texts: make(map[fontMeshCreator.FontType][]fontMeshCreator.GUIText)
	}
}

func (r *FontRenderer) Render() {
    r.Prepare()
    for font, texts := range r.texts {
        gl.ActiveTexture(gl.TEXTURE0)
        gl.BindTexture(gl.TEXTURE2D, font.TextureAtlas)
        for _, text := range texts {
            r.renderText(text)
        }
    }
}

func (r *FontRenderer) RemoveText(text fontMeshCreator.GUIText) {
    l := Texts[text.Font]
    l2 := make([]fontMeshCreator.GUIText, 0)
    for _, e := range l {
        if e != text {
            l2 = append(l2, e)
        }
    }
    if len(l2) > 0 {
        Texts[text.Font] = l2
    } else {
        delete(Texts, text.Font)
    }
}

func (r *FontRenderer) LoadText(text fontMeshCreator.GUIText) {
    font := text.Font
    data := font.LoadText(text)
    vao := LoadFontVAO(data.VertexesPositions, data.TexturesPositions)
    text.VaoID = vao
    text.VertexCount = data.VertexCount
    r.Texts[font] = append(r.Texts[font], text)
}

func (r *FontRenderer) CleanUp() {
	r.shader.CleanUp()
}

func (r *FontRenderer) Prepare() {
    gl.Enable(gl.BLEND)
    gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
    gl.Disable(gl.DEPTH_TEST)
    r.shader.Program.Start()
}

func (r *FontRenderer) renderText(text GUIText) {
    gl.BindVertexArray(text.Mesh)
    gl.EnableVertexAttribArray(0)
    gl.EnableVertexAttribArray(1)
    r.shader.LoadColour(text.Colour)
    r.shader.LoadTranslation(text.Position)
    gl.DrawArrays(gl.TRIANGLES, 0, text.VertexCount)
    gl.DisableVertexAttribArray(0)
    gl.DisableVertexAttribArray(1)
    gl.BindVertexArray(0)
}

func (r *FontRenderer) endRendering() {
    r.shader.Program.Stop()
    gl.Disable(gl.BLEND)
    gl.Enable(gl.DEPTH_TEST)
}
