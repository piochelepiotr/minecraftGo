package shaders

import (
	"github.com/go-gl/mathgl/mgl32"
)

type StaticShader struct {
    Program ShaderProgram
    TransformationMatrixLocation int32
    ProjectionMatrixLocation int32
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
    s.TransformationMatrixLocation = s.Program.GetUniformLocation("transformationMatrix\x00")
    s.ProjectionMatrixLocation = s.Program.GetUniformLocation("projectionMatrix\x00")
}

func (s *StaticShader) LoadTransformationMatrix(mat mgl32.Mat4) {
    s.Program.LoadMatrix4(s.TransformationMatrixLocation, mat)
}

func (s *StaticShader) LoadProjectionMatrix(mat mgl32.Mat4) {
    s.Program.LoadMatrix4(s.ProjectionMatrixLocation, mat)
}
