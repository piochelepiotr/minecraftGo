package renderEngine
import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/piochelepiotr/minecraftGo/models"
	texturesPackage "github.com/piochelepiotr/minecraftGo/textures"
	"image"
	"image/draw"
	_ "image/png"
    "os"
    "fmt"
)

var vaos = make([]uint32, 0)
var vbos = make([]uint32, 0)
var textures = make([]uint32, 0)

func LoadToVAO(positions []float32, textureCoord []float32, indices []uint32, normals []float32) models.RawModel {
    vaoID := createVAO()
    bindIndicesBuffer(indices)
    storeDataInAttributeList(0, 3, positions)
    storeDataInAttributeList(1, 2, textureCoord)
    storeDataInAttributeList(2, 3, normals)
    defer unbindVAO()
    return models.RawModel{
        VaoID:vaoID,
        VertexCount: int32(len(indices)),
    }
}

func loadTexture(file string) (uint32, error) {
	imgFile, err := os.Open(file)
	if err != nil {
		return 0, fmt.Errorf("texture %q not found on disk: %v", file, err)
	}
	img, _, err := image.Decode(imgFile)
	if err != nil {
		return 0, err
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		return 0, fmt.Errorf("unsupported stride")
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	var texture uint32
	gl.GenTextures(1, &texture)
    textures = append(textures, texture)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix))

	return texture, nil
}

func LoadTexture(file string) texturesPackage.ModelTexture {
    textureID, err := loadTexture(file)
    if err != nil {
        panic(err)
    }
    return texturesPackage.ModelTexture{
        Id:textureID,
        Reflectivity: 1,
        ShineDamper: 10,
        NumberOfRows: 1,
    }
}

func createVAO() uint32 {
    var vaoID uint32
    gl.GenVertexArrays(1, &vaoID)
    vaos = append(vaos, vaoID)
    gl.BindVertexArray(vaoID)
    return vaoID
}

func bindIndicesBuffer(indices []uint32) {
    var vboID uint32
    gl.GenBuffers(1, &vboID)
    vbos = append(vbos, vboID)
    gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, vboID)
    gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, 4*len(indices), gl.Ptr(indices), gl.STATIC_DRAW)
}


func storeDataInAttributeList(attributeNumber uint32, size int32, data []float32) {
    var vboID uint32
    gl.GenBuffers(1, &vboID)
    vbos = append(vbos, vboID)
    gl.BindBuffer(gl.ARRAY_BUFFER, vboID)
    gl.BufferData(gl.ARRAY_BUFFER, 4*len(data), gl.Ptr(data), gl.STATIC_DRAW)
    //gl.BufferSubData(gl.ARRAY_BUFFER, 0, 4*len(data), gl.Ptr(data))
    gl.VertexAttribPointer(attributeNumber, size, gl.FLOAT, false, 0, gl.PtrOffset(0))
    gl.BindBuffer(gl.ARRAY_BUFFER, 0)
}

func CleanUp() {
    for i := 0; i < len(vaos); i++ {
        gl.DeleteVertexArrays(1, &vaos[i])
    }
    for i := 0; i < len(vbos); i++ {
        gl.DeleteBuffers(1, &vbos[i])
    }
    for i := 0; i < len(textures); i++ {
        gl.DeleteTextures(1, &textures[i])
    }
}

func unbindVAO() {
    gl.BindVertexArray(0)
}

