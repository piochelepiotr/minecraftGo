package main
import (
    "fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
)
func LoadToVAO(positions []float64) RawModel {
    vaoID := createVAO()
    fmt.Println(vaoID)
    storeDataInAttributeList(0, positions)
    defer unbindVAO()
    return RawModel{}
}

func createVAO() int {
    vaoID := gl.GenVertexArrays()
    return vaoID
}

func storeDataInAttributeList(attributeNumber int, data []float64) {
}

func unbindVAO() {
    gl.BindVertexArray(0)
}

