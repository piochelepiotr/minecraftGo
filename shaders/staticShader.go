package shaders

import (
	"github.com/go-gl/mathgl/mgl32"
)

type StaticShader struct {
    Program ShaderProgram
    transformationMatrixLocation int32
    projectionMatrixLocation int32
    viewMatrixLocation int32
}

func CreateStaticShader() StaticShader {
    bindAttributes := func (s ShaderProgram) {
        s.bindAttribute(0, "position\x00")
        s.bindAttribute(1, "textureCoords\x00")
    }
    s:= StaticShader{Program: CreateShader("shaders/vertexShader.txt", "shaders/fragmentShader.txt", bindAttributes)}
    s.getAllUniformLocations()
    return s
}


func (s *StaticShader) getAllUniformLocations() {
    s.transformationMatrixLocation = s.Program.GetUniformLocation("transformationMatrix\x00")
    s.projectionMatrixLocation = s.Program.GetUniformLocation("projectionMatrix\x00")
    s.viewMatrixLocation = s.Program.GetUniformLocation("viewMatrix\x00")
}

func (s *StaticShader) LoadTransformationMatrix(mat mgl32.Mat4) {
    s.Program.LoadMatrix4(s.transformationMatrixLocation, mat)
}

func (s *StaticShader) LoadProjectionMatrix(mat mgl32.Mat4) {
    s.Program.LoadMatrix4(s.projectionMatrixLocation, mat)
}

func (s *StaticShader) LoadViewMatrix(mat mgl32.Mat4) {
    s.Program.LoadMatrix4(s.viewMatrixLocation, mat)
}
