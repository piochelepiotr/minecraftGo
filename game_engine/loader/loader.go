package loader

import (
	"fmt"
	"image"
	"image/draw"
	"time"

	// only support png images
	_ "image/png"

	"os"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/piochelepiotr/minecraftGo/game_engine/font"
	"github.com/piochelepiotr/minecraftGo/game_engine/guis"
	"github.com/piochelepiotr/minecraftGo/game_engine/models"
	"github.com/piochelepiotr/minecraftGo/game_engine/shaders"
	texturesPackage "github.com/piochelepiotr/minecraftGo/textures"
)

var vaos = make([]uint32, 0)
var vbos = make([]uint32, 0)
var textures = make([]uint32, 0)

var Debug bool

// LoadToVAO loads vertices into a vao
func LoadToVAO(positions []float32, textureCoord []float32, indices []uint32, normals []float32, colors []float32) models.RawModel {
	vaoID := createVAO()
	bindIndicesBuffer(indices)
	storeDataInAttributeList(0, 3, positions)
	storeDataInAttributeList(1, 2, textureCoord)
	storeDataInAttributeList(2, 3, normals)
	storeDataInAttributeList(3, 3, colors)
	defer unbindVAO()
	return models.RawModel{
		VaoID:       vaoID,
		VertexCount: int32(len(indices)),
	}
}

// LoadTexToVAO loads a texture (2D coords) into a VAO
func LoadTexToVAO(positions []float32) models.RawModel {
	vaoID := createVAO()
	storeDataInAttributeList(0, 2, positions)
	defer unbindVAO()
	return models.RawModel{
		VaoID:       vaoID,
		VertexCount: int32(len(positions) / 2),
	}
}

// LoadFontVAO loads a font into a VAO
func LoadFontVAO(positions []float32, textureCoord []float32) uint32 {
	vaoID := createVAO()
	storeDataInAttributeList(0, 2, positions)
	storeDataInAttributeList(1, 2, textureCoord)
	defer unbindVAO()
	return vaoID
}

// LoadFont create a font
func LoadFont(fontTexture, fontFile string, aspectRatio float32) *font.FontType {
	textureID, err := loadTexture(fontTexture)
	if err != nil {
		panic(err)
	}
	return font.CreateFontType(textureID, fontFile, aspectRatio)

}

// CreateGuiRenderer returns a gui renderer
func CreateGuiRenderer() guis.GuiRenderer {
	positions := []float32{
		-1, 1,
		-1, -1,
		1, 1,
		1, -1,
	}
	return guis.GuiRenderer{
		Quad:   LoadTexToVAO(positions),
		Shader: shaders.CreateGuiShader(),
	}
}

// LoadText loads a text into a VAO
func LoadText(text font.GUIText) font.GUIText {
	font := text.Font
	data := font.LoadText(text)
	vao := LoadFontVAO(data.VertexPositions, data.TextureCoords)
	text.TextMeshVao = vao
	text.VertexCount = data.GetVertexCount()
	return text
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
	gl.GenerateMipmap(gl.TEXTURE_2D)
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

// LoadModelTexture gets a texture from a file and loads it
func LoadModelTexture(file string, numberOfRows uint32) texturesPackage.ModelTexture {
	textureID, err := loadTexture(file)
	if err != nil {
		panic(err)
	}
	return texturesPackage.ModelTexture{
		Id:           textureID,
		NumberOfRows: numberOfRows,
	}
}

//LoadGuiTexture loads a texture into a VAO
func LoadGuiTexture(file string, position, scale mgl32.Vec2) guis.GuiTexture {
	textureID, err := loadTexture(file)
	if err != nil {
		panic(err)
	}
	return guis.GuiTexture{
		Id:       textureID,
		Position: position,
		Scale:    scale,
	}
}

func createVAO() uint32 {
	if Debug {
		fmt.Println("hello")
		time.Sleep(time.Second)
	}
	var vaoID uint32
	gl.GenVertexArrays(1, &vaoID)
	if Debug {
		fmt.Println("after gen")
		time.Sleep(time.Second)
	}
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
	gl.VertexAttribPointer(attributeNumber, size, gl.FLOAT, false, 0, gl.PtrOffset(0))
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
}

// CleanUp clears all VAOs and VBOs
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
