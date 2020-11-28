package loader

import (
	"fmt"
	"image"
	"image/draw"
	// only support png images
	_ "image/png"

	"os"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/piochelepiotr/minecraftGo/game_engine/font"
	"github.com/piochelepiotr/minecraftGo/game_engine/guis"
	"github.com/piochelepiotr/minecraftGo/game_engine/models"
	"github.com/piochelepiotr/minecraftGo/game_engine/shaders"
)

type vao struct {
	id uint32
	vbo []uint32
}

type Loader struct {
	vaos []uint32
	vbos []uint32
	textures []uint32
}

func NewLoader() *Loader {
	return &Loader{
		vaos : make([]uint32, 0),
		vbos : make([]uint32, 0),
		textures : make([]uint32, 0),

	}
}

// LoadToVAO loads vertices into a vao
func (l *Loader) LoadToVAO(positions []float32, textureCoord []float32, indices []uint32, normals []float32, colors []float32) models.RawModel {
	vaoID := l.createVAO()
	l.bindIndicesBuffer(indices)
	l.storeDataInAttributeList(0, 3, positions)
	l.storeDataInAttributeList(1, 2, textureCoord)
	l.storeDataInAttributeList(2, 3, normals)
	l.storeDataInAttributeList(3, 3, colors)
	defer l.unbindVAO()
	return models.RawModel{
		VaoID:       vaoID,
		VertexCount: int32(len(indices)),
	}
}

func (l *Loader) LoadLinesToVAO(positions []float32, indices []uint32) models.RawModel{
	vaoID := l.createVAO()
	l.bindIndicesBuffer(indices)
	l.storeDataInAttributeList(0, 3, positions)
	defer l.unbindVAO()
	return models.RawModel{
		VaoID:       vaoID,
		VertexCount: int32(len(indices)),
	}
}

// LoadTexToVAO loads a texture (2D coords) into a VAO
func (l *Loader) LoadTexToVAO(positions []float32) models.RawModel {
	vaoID := l.createVAO()
	l.storeDataInAttributeList(0, 2, positions)
	defer l.unbindVAO()
	return models.RawModel{
		VaoID:       vaoID,
		VertexCount: int32(len(positions) / 2),
	}
}

// LoadFontVAO loads a font into a VAO
func (l *Loader) LoadFontVAO(positions []float32, textureCoord []float32) uint32 {
	vaoID := l.createVAO()
	l.storeDataInAttributeList(0, 2, positions)
	l.storeDataInAttributeList(1, 2, textureCoord)
	defer l.unbindVAO()
	return vaoID
}

// LoadFont create a font
func (l *Loader) LoadFont(fontTexture, fontFile string, aspectRatio float32) *font.FontType {
	textureID, err := l.loadTexture(fontTexture)
	if err != nil {
		panic(err)
	}
	return font.CreateFontType(textureID, fontFile, aspectRatio)

}

// CreateGuiRenderer returns a gui renderer
func (l *Loader) CreateGuiRenderer() guis.GuiRenderer {
	positions := []float32{
		-1, 1,
		-1, -1,
		1, 1,
		1, -1,
	}
	return guis.GuiRenderer{
		Quad:   l.LoadTexToVAO(positions),
		Shader: shaders.CreateGuiShader(),
	}
}

// LoadText loads a text into a VAO
func (l *Loader) LoadText(text font.GUIText) font.GUIText {
	font := text.Font
	data := font.LoadText(text)
	vao := l.LoadFontVAO(data.VertexPositions, data.TextureCoords)
	text.TextMeshVao = vao
	text.VertexCount = data.GetVertexCount()
	text.Width = data.Width
	return text
}

func (l *Loader) loadTexture(file string) (uint32, error) {
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
	l.textures = append(l.textures, texture)
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
func (l *Loader) LoadModelTexture(file string) uint32 {
	textureID, err := l.loadTexture(file)
	if err != nil {
		panic(err)
	}
	return textureID
}

//LoadGuiTexture loads a texture into a VAO
func (l *Loader) LoadGuiTexture(file string, position, scale mgl32.Vec2) guis.GuiTexture {
	textureID, err := l.loadTexture(file)
	if err != nil {
		panic(err)
	}
	return guis.GuiTexture{
		Id:       textureID,
		Position: position,
		Scale:    scale,
	}
}

func (l *Loader) createVAO() uint32 {
	var vaoID uint32
	gl.GenVertexArrays(1, &vaoID)
	l.vaos = append(l.vaos, vaoID)
	gl.BindVertexArray(vaoID)
	return vaoID
}

func (l *Loader) bindIndicesBuffer(indices []uint32) {
	var vboID uint32
	gl.GenBuffers(1, &vboID)
	l.vbos = append(l.vbos, vboID)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, vboID)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, 4*len(indices), gl.Ptr(indices), gl.STATIC_DRAW)
}

func (l *Loader) storeDataInAttributeList(attributeNumber uint32, size int32, data []float32) {
	var vboID uint32
	gl.GenBuffers(1, &vboID)
	l.vbos = append(l.vbos, vboID)
	gl.BindBuffer(gl.ARRAY_BUFFER, vboID)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(data), gl.Ptr(data), gl.STATIC_DRAW)
	gl.VertexAttribPointer(attributeNumber, size, gl.FLOAT, false, 0, gl.PtrOffset(0))
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
}

// CleanUp clears all VAOs and VBOs
func (l *Loader) CleanUp() {
	for i := 0; i < len(l.vaos); i++ {
		gl.DeleteVertexArrays(1, &l.vaos[i])
	}
	for i := 0; i < len(l.vbos); i++ {
		gl.DeleteBuffers(1, &l.vbos[i])
	}
	for i := 0; i < len(l.textures); i++ {
		gl.DeleteTextures(1, &l.textures[i])
	}
}

func (l *Loader) unbindVAO() {
	gl.BindVertexArray(0)
}