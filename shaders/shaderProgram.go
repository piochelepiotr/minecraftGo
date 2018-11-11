package shaders
import (
    "fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
    "os"
    "io/ioutil"
    "strings"
)

type ShaderProgram struct {
    ProgramID uint32
    VertexShaderId uint32
    FragmentShaderId uint32
}

func CreateShader(vertexShader string, fragmentShader string) ShaderProgram {
    s := ShaderProgram{}
    s.VertexShaderId = loadShader(vertexShader, gl.VERTEX_SHADER)
    s.FragmentShaderId = loadShader(fragmentShader, gl.FRAGMENT_SHADER)
    s.ProgramID = gl.CreateProgram()
    gl.AttachShader(s.ProgramID, s.VertexShaderId)
    gl.AttachShader(s.ProgramID, s.FragmentShaderId)
    gl.LinkProgram(s.ProgramID)
    gl.ValidateProgram(s.ProgramID)
    return s
}

func (s ShaderProgram) Start() {
    gl.UseProgram(s.ProgramID)
}

func (s ShaderProgram) Stop() {
    gl.UseProgram(0)
}

func (s ShaderProgram) CleanUp() {
    s.Stop()
    gl.DetachShader(s.ProgramID, s.VertexShaderId)
    gl.DetachShader(s.ProgramID, s.FragmentShaderId)
    gl.DeleteShader(s.VertexShaderId)
    gl.DeleteShader(s.FragmentShaderId)
    gl.DeleteProgram(s.ProgramID)
}

func (s ShaderProgram) bindAttribute(attribute uint32, variableName string) {
    gl.BindAttribLocation(s.ProgramID, attribute, gl.Str(variableName))
}

func loadShader (file string, shaderType uint32) uint32 {
    dat, err := ioutil.ReadFile(file)
    if err != nil {
        panic(err)
    }
    shaderID := gl.CreateShader(shaderType)
    if shaderID == 0 {
        fmt.Println("Shader type doesn't exist")
        os.Exit(1)
    }
	csources, free := gl.Strs(string(dat) + "\x00")
	gl.ShaderSource(shaderID, 1, csources, nil)
	free()
	gl.CompileShader(shaderID)
	var status int32
	gl.GetShaderiv(shaderID, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
        fmt.Println("Error while compiling shader ", file)
		var logLength int32
		gl.GetShaderiv(shaderID, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shaderID, logLength, nil, gl.Str(log))
        fmt.Println(log)
		return 0
	}
    return shaderID
}
