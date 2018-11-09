package main
import (
    "fmt"
	"github.com/go-gl/gl/v4.6-core/gl"
    "unsafe"
)

var vaos = make([]uint32, 0)
var vbos = make([]uint32, 0)

func LoadToVAO(positions []float64) RawModel {
    vaoID := createVAO()
    fmt.Println(vaoID)
    storeDataInAttributeList(0, positions)
    defer unbindVAO()
    return RawModel{}
}

func createVAO() uint32 {
    var vaoID uint32
    gl.GenVertexArrays(1, &vaoID)
    append(vaos, vaoID)
    gl.BindVertexArray(vaoID)
    return vaoID
}

func storeDataInAttributeList(attributeNumber uint32, data []float64) {
    var vboID uint32
    gl.GenBuffers(1, &vboID)
    append(vbos, vboID)
    gl.BindBuffer(vboID, gl.ARRAY_BUFFER)
    gl.BufferData(vboID, unsafe.Sizeof(float64)*len(data), gl.Ptr(data), gl.STATIC_DRAW)
    gl.VertexAttribPointer(attributeNumber, 3, gl.FLOAT, false, 0, 0)
    gl.BindBuffer(0, gl.ARRAY_BUFFER)
}

func cleanUp() {
    for i, vaoID := range vaos {
        gl.DeleteVertexArrays(
    }
    for i, vboID := range vbos {
        gl.DeleteBuffers(
    }
}

func unbindVAO() {
    gl.BindVertexArray(0)
}

