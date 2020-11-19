package shaders

import (
	"github.com/go-gl/mathgl/mgl32"
    "github.com/piochelepiotr/minecraftGo/entities"
    "github.com/piochelepiotr/minecraftGo/toolbox"
)

type OutlineShader struct {
    Program ShaderProgram
    transformationMatrixLocation int32
    projectionMatrixLocation int32
    viewMatrixLocation int32
}

func CreateOutlineShader() OutlineShader {
    bindAttributes := func (s ShaderProgram) {
        s.bindAttribute(0, "position\x00")
    }
    s:= OutlineShader{Program: CreateShader("outlineShader.vert", "outlineShader.frag", bindAttributes)}
    s.getAllUniformLocations()
    return s
}


func (s *OutlineShader) getAllUniformLocations() {
    s.transformationMatrixLocation = s.Program.GetUniformLocation("transformationMatrix\x00")
    s.projectionMatrixLocation = s.Program.GetUniformLocation("projectionMatrix\x00")
    s.viewMatrixLocation = s.Program.GetUniformLocation("viewMatrix\x00")
}

func (s *OutlineShader) LoadTransformationMatrix(mat mgl32.Mat4) {
    s.Program.LoadMatrix4(s.transformationMatrixLocation, mat)
}

func (s *OutlineShader) LoadProjectionMatrix(mat mgl32.Mat4) {
    s.Program.LoadMatrix4(s.projectionMatrixLocation, mat)
}

func (s *OutlineShader) LoadViewMatrix(camera *entities.Camera) {
    viewMatrix := toolbox.CreateViewMatrix(camera.Position, camera.Rotation)
    s.Program.LoadMatrix4(s.viewMatrixLocation, viewMatrix)
}

func (s *OutlineShader) CleanUp() {
    s.Program.CleanUp()
}
