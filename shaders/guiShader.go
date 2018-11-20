package shaders

import (
	"github.com/go-gl/mathgl/mgl32"
)

type GuiShader struct {
    Program ShaderProgram
    transformationMatrixLocation int32
}

func CreateGuiShader() GuiShader {
    bindAttributes := func (s ShaderProgram) {
        s.bindAttribute(0, "position\x00")
    }
    s:= GuiShader{Program: CreateShader("shaders/guiShader.vert", "shaders/guiShader.frag", bindAttributes)}
    s.getAllUniformLocations()
    return s
}


func (s *GuiShader) getAllUniformLocations() {
    s.transformationMatrixLocation = s.Program.GetUniformLocation("transformationMatrix\x00")
}

func (s *GuiShader) LoadTransformationMatrix(mat mgl32.Mat4) {
    s.Program.LoadMatrix4(s.transformationMatrixLocation, mat)
}

func (s *GuiShader) CleanUp() {
    s.Program.CleanUp()
}
