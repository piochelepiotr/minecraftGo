package renderEngine


import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/piochelepiotr/minecraftGo/models"
    "io/ioutil"
    "strings"
    "strconv"
)

func toFloat(s string) float32 {
    value, err := strconv.ParseFloat(s, 32)
    if err != nil {
        panic(err)
    }
    return float32(value)
}

func toInt(s string) uint32 {
    value, err := strconv.ParseUint(s, 10, 32)
    if err != nil {
        panic(err)
    }
    return uint32(value)
}


func LoadObjModel(file string) models.RawModel {
    dat, err := ioutil.ReadFile(file)
    if err != nil {
        panic(err)
    }
    vertices := make([]mgl32.Vec3, 0)
    textures := make([]mgl32.Vec2, 0)
    normals := make([]mgl32.Vec3, 0)
    var verticesArray []float32
    var texturesArray []float32
    var normalsArray []float32
    indicesArray := make([]uint32, 0)
    lines := strings.Split(string(dat), "\n")
    var fStart int
    for i, line := range lines {
        splited := strings.Split(line, " ")
        if len(splited) == 0 {
            continue
        }
        t := splited[0]
        if t == "v" {
            vertices = append(vertices, mgl32.Vec3{toFloat(splited[1]),toFloat(splited[2]),toFloat(splited[3])})
        }
        if t == "vn" {
            normals = append(normals, mgl32.Vec3{toFloat(splited[1]),toFloat(splited[2]),toFloat(splited[3])})
        }
        if t == "vt" {
            textures = append(textures, mgl32.Vec2{toFloat(splited[1]),toFloat(splited[2])})
        }
        if t == "f" {
            fStart = i
            texturesArray = make([]float32, len(vertices)*2)
            normalsArray = make([]float32, len(vertices)*3)
            verticesArray = make([]float32, len(vertices)*3)
            break
        }
    }
    for i := fStart; i < len(lines); i++ {
        splited := strings.Split(lines[i], " ")
        if len(splited) == 0 || splited[0] != "f" {
            break
        }
        vertex1 := strings.Split(splited[1], "/")
        vertex2 := strings.Split(splited[2], "/")
        vertex3 := strings.Split(splited[3], "/")
        indicesArray = processVertex(vertex1, indicesArray, textures, normals, vertices, texturesArray, normalsArray, verticesArray)
        indicesArray = processVertex(vertex2, indicesArray, textures, normals, vertices, texturesArray, normalsArray, verticesArray)
        indicesArray = processVertex(vertex3, indicesArray, textures, normals, vertices, texturesArray, normalsArray, verticesArray)
    }
    return LoadToVAO(verticesArray, texturesArray, indicesArray, normalsArray)
}

func processVertex(vertexData []string, indices []uint32, textures []mgl32.Vec2, normals []mgl32.Vec3,vertices []mgl32.Vec3, texturesArray []float32, normalsArray[]float32, verticesArray []float32) []uint32 {
    currentVertexPointer := toInt(vertexData[0]) - 1
    indices = append(indices, currentVertexPointer)
    vert := vertices[currentVertexPointer]
    verticesArray[3*currentVertexPointer] = vert.X()
    verticesArray[3*currentVertexPointer+1] = vert.Y()
    verticesArray[3*currentVertexPointer+2] = vert.Z()
    tex := textures[toInt(vertexData[1])-1]
    texturesArray[2*currentVertexPointer] = tex.X()
    texturesArray[2*currentVertexPointer+1] = 1- tex.Y()
    norm := normals[toInt(vertexData[2])-1]
    normalsArray[3*currentVertexPointer] = norm.X()
    normalsArray[3*currentVertexPointer+1] = norm.Y()
    normalsArray[3*currentVertexPointer+2] = norm.Z()
    return indices
}
