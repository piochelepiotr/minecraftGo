package render

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/piochelepiotr/minecraftGo/game_engine/font"
	"github.com/piochelepiotr/minecraftGo/game_engine/shaders"
)

type FontRenderer struct {
	shader shaders.FontShader
	texts  map[*font.FontType][]font.GUIText
}

func CreateFontRenderer() FontRenderer {
	return FontRenderer{
		shader: shaders.CreateFontShader(),
		texts:  make(map[*font.FontType][]font.GUIText),
	}
}

func (r *FontRenderer) Render(texts []font.GUIText) {
	r.LoadTexts(texts)
	r.Prepare()
	for font, texts := range r.texts {
		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, font.TextureAtlas)
		for _, text := range texts {
			r.renderText(text)
		}
	}
	r.endRendering()
}

func (r *FontRenderer) LoadText(text font.GUIText) {
	font := text.Font
	r.texts[font] = append(r.texts[font], text)
}

func (r *FontRenderer) LoadTexts(texts []font.GUIText) {
	for _, text := range texts {
		r.LoadText(text)
	}
}

//func (r *FontRenderer) RemoveText(text font.GUIText) {
//	l := r.texts[text.Font]
//	l2 := make([]font.GUIText, 0)
//	for _, e := range l {
//		if e != text {
//			l2 = append(l2, e)
//		}
//	}
//	if len(l2) > 0 {
//		r.texts[text.Font] = l2
//	} else {
//		delete(r.texts, text.Font)
//	}
//}

func (r *FontRenderer) CleanUp() {
	r.shader.CleanUp()
}

func (r *FontRenderer) Prepare() {

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.Disable(gl.DEPTH_TEST)
	r.shader.Program.Start()
}

func (r *FontRenderer) renderText(text font.GUIText) {
	gl.BindVertexArray(text.TextMeshVao)
	gl.EnableVertexAttribArray(0)
	gl.EnableVertexAttribArray(1)
	r.shader.LoadColour(text.Colour)
	//r.shader.LoadTranslation(text.Position)
	r.shader.LoadTranslation(text.Position.Add(mgl32.Vec2{0, -text.GetLineHeight() / 2}))
	gl.DrawArrays(gl.TRIANGLES, 0, text.VertexCount)
	gl.DisableVertexAttribArray(0)
	gl.DisableVertexAttribArray(1)
	gl.BindVertexArray(0)
}

func (r *FontRenderer) endRendering() {
	r.shader.Program.Stop()
	gl.Disable(gl.BLEND)
	gl.Enable(gl.DEPTH_TEST)
	r.texts = make(map[*font.FontType][]font.GUIText)
}
