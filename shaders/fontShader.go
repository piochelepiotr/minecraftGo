package shaders

import (
	"github.com/go-gl/mathgl/mgl32"
)

type FontShader struct {
    Program ShaderProgram
    ColourLocation int32
    TranslationLocation int32
}

func CreateFontShader() FontShader {
    bindAttributes := func (s ShaderProgram) {
        s.bindAttribute(0, "position\x00")
        s.bindAttribute(1, "textureCoords\x00")
    }
    s:= GuiShader{Program: CreateShader("shaders/fontShader.vert", "shaders/fontShader.frag", bindAttributes)}
    s.getAllUniformLocations()
    return s
}

func (s *FontShader) getAllUniformLocations() {
    s.ColourLocation = s.Program.GetUniformLocation("colour\x00")
    s.TranslationLocation = s.Program.GetUniformLocation("translation\x00")
}

func (s *FontShader) LoadColour(colour mgl32.Vec3) {
    s.Program.LoadVector(s.ColourLocation, colour)
}

func (s *FontShader) LoadTranslation(translation mgl32.Vec2) {
    s.Program.Load2DVector(s.TranslationLocation, translation)
}

func (s *FontShader) CleanUp() {
    s.Program.CleanUp()
}
