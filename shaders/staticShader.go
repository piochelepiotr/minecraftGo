package shaders

import (
	"github.com/go-gl/mathgl/mgl32"
    "github.com/piochelepiotr/minecraftGo/entities"
    "github.com/piochelepiotr/minecraftGo/toolbox"
)

type StaticShader struct {
    Program ShaderProgram
    transformationMatrixLocation int32
    projectionMatrixLocation int32
    viewMatrixLocation int32
    lightPositionLocation int32
    lightColourLocation int32
    reflectivityLocation int32
    shineDamperLocation int32
}

func CreateStaticShader() StaticShader {
    bindAttributes := func (s ShaderProgram) {
        s.bindAttribute(0, "position\x00")
        s.bindAttribute(1, "textureCoords\x00")
        s.bindAttribute(2, "normal\x00")
        s.bindAttribute(3, "colors\x00")
    }
    s:= StaticShader{Program: CreateShader("shaders/vertexShader.txt", "shaders/fragmentShader.txt", bindAttributes)}
    s.getAllUniformLocations()
    return s
}


func (s *StaticShader) getAllUniformLocations() {
    s.transformationMatrixLocation = s.Program.GetUniformLocation("transformationMatrix\x00")
    s.projectionMatrixLocation = s.Program.GetUniformLocation("projectionMatrix\x00")
    s.viewMatrixLocation = s.Program.GetUniformLocation("viewMatrix\x00")
    s.lightPositionLocation = s.Program.GetUniformLocation("lightPosition\x00")
    s.lightColourLocation = s.Program.GetUniformLocation("lightColour\x00")
    s.reflectivityLocation = s.Program.GetUniformLocation("reflectivity\x00")
    s.shineDamperLocation = s.Program.GetUniformLocation("shineDamper\x00")
}

func (s *StaticShader) LoadTransformationMatrix(mat mgl32.Mat4) {
    s.Program.LoadMatrix4(s.transformationMatrixLocation, mat)
}

func (s *StaticShader) LoadProjectionMatrix(mat mgl32.Mat4) {
    s.Program.LoadMatrix4(s.projectionMatrixLocation, mat)
}

func (s *StaticShader) LoadViewMatrix(camera *entities.Camera) {
    viewMatrix := toolbox.CreateViewMatrix(camera.Position, camera.Rotation)
    s.Program.LoadMatrix4(s.viewMatrixLocation, viewMatrix)
}

func (s *StaticShader) LoadLight(light *entities.Light) {
    s.Program.LoadVector(s.lightPositionLocation, light.Position)
    s.Program.LoadVector(s.lightColourLocation, light.Colour)
}

func (s *StaticShader) LoadShineVariables(shineDamper float32, reflectivity float32) {
    s.Program.LoadFloat(s.shineDamperLocation, shineDamper)
    s.Program.LoadFloat(s.reflectivityLocation, reflectivity)
}

func (s *StaticShader) CleanUp() {
    s.Program.CleanUp()
}
