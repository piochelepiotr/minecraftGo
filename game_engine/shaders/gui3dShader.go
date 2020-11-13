package shaders

import (
	"github.com/go-gl/mathgl/mgl32"
    "github.com/piochelepiotr/minecraftGo/entities"
    "github.com/piochelepiotr/minecraftGo/toolbox"
)

type Gui3dShader struct {
    Program ShaderProgram
    transformationMatrixLocation int32
    projectionMatrixLocation int32
    viewMatrixLocation int32
    aspectRatioLocation int32
    translationLocation int32
}

func CreateGui3dShader() Gui3dShader {
    bindAttributes := func (s ShaderProgram) {
        s.bindAttribute(0, "position\x00")
        s.bindAttribute(1, "textureCoords\x00")
        s.bindAttribute(2, "normal\x00")
        s.bindAttribute(3, "colors\x00")
    }
    s:= Gui3dShader{Program: CreateShader("gui3dShader.vert", "fragmentShader.txt", bindAttributes)}
    s.getAllUniformLocations()
    return s
}


func (s *Gui3dShader) getAllUniformLocations() {
    s.transformationMatrixLocation = s.Program.GetUniformLocation("transformationMatrix\x00")
    s.projectionMatrixLocation = s.Program.GetUniformLocation("projectionMatrix\x00")
    s.viewMatrixLocation = s.Program.GetUniformLocation("viewMatrix\x00")
    s.aspectRatioLocation = s.Program.GetUniformLocation("aspectRatio\x00")
    s.translationLocation = s.Program.GetUniformLocation("translation\x00")
}

func (s *Gui3dShader) LoadTransformationMatrix(mat mgl32.Mat4) {
    s.Program.LoadMatrix4(s.transformationMatrixLocation, mat)
}

func (s *Gui3dShader) LoadProjectionMatrix(mat mgl32.Mat4) {
    s.Program.LoadMatrix4(s.projectionMatrixLocation, mat)
}

func (s *Gui3dShader) LoadViewMatrix(camera *entities.Camera) {
    viewMatrix := toolbox.CreateViewMatrix(camera.Position, camera.Rotation)
    s.Program.LoadMatrix4(s.viewMatrixLocation, viewMatrix)
}

func (s *Gui3dShader) LoadAspectRatio(aspectRatio float32) {
    s.Program.LoadFloat(s.aspectRatioLocation, aspectRatio)
}

func (s *Gui3dShader) LoadTranslation(translation mgl32.Vec2) {
    s.Program.Load2DVector(s.translationLocation, translation)
}

func (s *Gui3dShader) CleanUp() {
    s.Program.CleanUp()
}
