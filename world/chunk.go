package world

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/piochelepiotr/minecraftGo/loader"
	"github.com/piochelepiotr/minecraftGo/models"
	"github.com/piochelepiotr/minecraftGo/textures"
)

// Chunk is set cube of blocks
type Chunk struct {
	Model  models.RawModel
	blocks []Block
	generator *Generator
	Start  Point
}

// NumberRowsTextures is the number number of rows on the texture image
const NumberRowsTextures int = 16

// ChunkSize is the size of a chunk in blocks
const ChunkSize int = 16

// ChunkSize2 is the area of a chunk in blocks
const ChunkSize2 = ChunkSize * ChunkSize

// ChunkSize3 is the volume of a chunk in blocks
const ChunkSize3 = ChunkSize2 * ChunkSize


// CreateChunk allows you to create a chunk by passing the start point (the second chunk is at position ChunkSize-1)
func CreateChunk(startX int, startY int, startZ int, modelTexture textures.ModelTexture, generator *Generator) Chunk {
	var chunk Chunk
	chunk.Start = Point{
		X: startX,
		Y: startY,
		Z: startZ,
	}
	chunk.blocks = make([]Block, ChunkSize3)
	chunk.generator = generator
	for x := 0; x < ChunkSize; x++ {
		for z := 0; z < ChunkSize; z++ {
			for y := 0; y < ChunkSize; y++ {
				chunk.setBlockNoUpdate(x, y, z, chunk.generator.BlockType(startX + x, startY + y, startZ + z))
			}
		}
	}
	chunk.buildFaces()
	return chunk
}

// setBlockNoUpdate sets a block in a chunk, it doesn't refresh the display
func (c *Chunk) setBlockNoUpdate(x, y, z int, b Block) {
	c.blocks[x*ChunkSize2+y*ChunkSize+z] = b
}

// SetBlock sets a block in a chunk and refreshes the model
func (c *Chunk) SetBlock(x, y, z int, b Block) {
	c.blocks[x*ChunkSize2+y*ChunkSize+z] = b
	c.buildFaces()
	//fmt.Println("hello")
	//c.Model.VertexCount = 0
}

// GetBlock gets the block of a chunk
func (c *Chunk) GetBlock(x, y, z int) Block {
	return c.blocks[x*ChunkSize2+y*ChunkSize+z]
}

// GetHeight gets the height of the chunk in blocks (not including air)
func (c *Chunk) GetHeight(x, z int) int {
	for y := ChunkSize - 1; y >= 0; y-- {
		if c.GetBlock(x, y, z) != Air {
			return y + 1
		}
	}
	return 0
}

func (c *Chunk) buildFaces() {
	vertices := make([]mgl32.Vec3, 0)
	textures := make([]mgl32.Vec2, 0)
	normals := make([]mgl32.Vec3, 0)
	colors := make([]mgl32.Vec3, 0)
	indexes := make([]uint32, 0)
	nextIndex := uint32(0)
	for x := 0; x < ChunkSize; x++ {
		for y := 0; y < ChunkSize; y++ {
			for z := 0; z < ChunkSize; z++ {
				b := c.GetBlock(x, y, z)
				a := Air
				//add face if not block isn't air and the block next to it is air
				/*faces are :
				 * up ( + y)
				 * bottom (-y)
				 * right (+x)
				 * left (-x)
				 * front (+z)
				 * back (-z)
				 * 0 point is : bottom, left, back
				 */
				//up
				xF := float32(x)
				yF := float32(y)
				zF := float32(z)
				if !(b == a || (y+1 < ChunkSize && c.GetBlock(x, y+1, z) != a)) {
					n := mgl32.Vec3{0, 1, 0}
					p1 := mgl32.Vec3{xF, yF + 1, zF}
					p2 := mgl32.Vec3{xF + 1, yF + 1, zF}
					p3 := mgl32.Vec3{xF + 1, yF + 1, zF + 1}
					p4 := mgl32.Vec3{xF, yF + 1, zF + 1}
					nextIndex = c.addFace(&vertices, &textures, &normals, &colors, &indexes, p1, p2, p3, p4, n, b.GetSide(Top), nextIndex, true)
				}
				//bottom
				if !(b == a || (y-1 > 0 && c.GetBlock(x, y-1, z) != a)) {
					n := mgl32.Vec3{0, -1, 0}
					p1 := mgl32.Vec3{xF, yF, (zF)}
					p2 := mgl32.Vec3{xF + 1, yF, (zF)}
					p3 := mgl32.Vec3{xF + 1, yF, (zF + 1)}
					p4 := mgl32.Vec3{xF, yF, (zF + 1)}
					nextIndex = c.addFace(&vertices, &textures, &normals, &colors, &indexes, p1, p2, p3, p4, n, b.GetSide(Bottom), nextIndex, false)
				}
				//right
				if !(b == a || (x+1 < ChunkSize && c.GetBlock(x+1, y, z) != a)) {
					n := mgl32.Vec3{1, 0, 0}
					p1 := mgl32.Vec3{(xF + 1), (yF + 1), (zF + 1)}
					p2 := mgl32.Vec3{(xF + 1), (yF + 1), (zF)}
					p3 := mgl32.Vec3{(xF + 1), (yF), (zF)}
					p4 := mgl32.Vec3{(xF + 1), (yF), (zF + 1)}
					nextIndex = c.addFace(&vertices, &textures, &normals, &colors, &indexes, p1, p2, p3, p4, n, b.GetSide(Side), nextIndex, true)
				}
				//left
				if !(b == a || (x-1 > 0 && c.GetBlock(x-1, y, z) != a)) {
					n := mgl32.Vec3{-1, 0, 0}
					p1 := mgl32.Vec3{(xF), (yF + 1), (zF)}
					p2 := mgl32.Vec3{(xF), (yF + 1), (zF + 1)}
					p3 := mgl32.Vec3{(xF), (yF), (zF + 1)}
					p4 := mgl32.Vec3{(xF), (yF), (zF)}
					nextIndex = c.addFace(&vertices, &textures, &normals, &colors, &indexes, p1, p2, p3, p4, n, b.GetSide(Side), nextIndex, true)
				}
				//front
				if !(b == a || (z+1 < ChunkSize && c.GetBlock(x, y, z+1) != a)) {
					n := mgl32.Vec3{0, 0, 1}
					p1 := mgl32.Vec3{(xF), (yF + 1), (zF + 1)}
					p2 := mgl32.Vec3{(xF + 1), (yF + 1), (zF + 1)}
					p3 := mgl32.Vec3{(xF + 1), (yF), (zF + 1)}
					p4 := mgl32.Vec3{(xF), (yF), (zF + 1)}
					nextIndex = c.addFace(&vertices, &textures, &normals, &colors, &indexes, p1, p2, p3, p4, n, b.GetSide(Side), nextIndex, true)
				}
				//back
				if !(b == a || (z-1 > 0 && c.GetBlock(x, y, z-1) != a)) {
					n := mgl32.Vec3{0, 0, -1}
					p1 := mgl32.Vec3{(xF + 1), (yF + 1), (zF)}
					p2 := mgl32.Vec3{(xF), (yF + 1), (zF)}
					p3 := mgl32.Vec3{(xF), (yF), (zF)}
					p4 := mgl32.Vec3{(xF + 1), (yF), (zF)}
					nextIndex = c.addFace(&vertices, &textures, &normals, &colors, &indexes, p1, p2, p3, p4, n, b.GetSide(Side), nextIndex, true)
				}
			}
		}
	}
	c.buildRawModel(vertices, textures, normals, colors, indexes)
}

func (c *Chunk) addFace(vertices *[]mgl32.Vec3, textures *[]mgl32.Vec2, normals *[]mgl32.Vec3, colors *[]mgl32.Vec3, indexes *[]uint32, p1, p2, p3, p4, n mgl32.Vec3, b Block, nextIndex uint32, inverseRotation bool) uint32 {
	if b == Grass {
		if n.ApproxEqual(mgl32.Vec3{0, -1, 0}) {
			b = Dirt
		} else if !n.ApproxEqual(mgl32.Vec3{0, 1, 0}) {
			b = Grass
		}
	}
	textureX := int(b) % NumberRowsTextures
	textureY := int(b) / NumberRowsTextures
	offsetTextureX := float32(textureX) / float32(NumberRowsTextures)
	offsetTextureY := float32(textureY) / float32(NumberRowsTextures)
	offsetTexture := mgl32.Vec2{offsetTextureX, offsetTextureY}
	t1 := mgl32.Vec2{0, 0}
	t2 := mgl32.Vec2{1, 0}
	t3 := mgl32.Vec2{1, 1}
	t4 := mgl32.Vec2{0, 1}
	t1 = t1.Mul(1.0 / float32(NumberRowsTextures)).Add(offsetTexture)
	t2 = t2.Mul(1.0 / float32(NumberRowsTextures)).Add(offsetTexture)
	t3 = t3.Mul(1.0 / float32(NumberRowsTextures)).Add(offsetTexture)
	t4 = t4.Mul(1.0 / float32(NumberRowsTextures)).Add(offsetTexture)
	*vertices = append(*vertices, p1)
	*vertices = append(*vertices, p2)
	*vertices = append(*vertices, p3)
	*vertices = append(*vertices, p4)
	*normals = append(*normals, n)
	*normals = append(*normals, n)
	*normals = append(*normals, n)
	*normals = append(*normals, n)

	color := b.GetColor()
	*colors = append(*colors, color)
	*colors = append(*colors, color)
	*colors = append(*colors, color)
	*colors = append(*colors, color)
	*textures = append(*textures, t1)
	*textures = append(*textures, t2)
	*textures = append(*textures, t3)
	*textures = append(*textures, t4)
	if inverseRotation {
		*indexes = append(*indexes, nextIndex)
		*indexes = append(*indexes, nextIndex+2)
		*indexes = append(*indexes, nextIndex+1)

		*indexes = append(*indexes, nextIndex)
		*indexes = append(*indexes, nextIndex+3)
		*indexes = append(*indexes, nextIndex+2)
	} else {
		*indexes = append(*indexes, nextIndex)
		*indexes = append(*indexes, nextIndex+1)
		*indexes = append(*indexes, nextIndex+2)

		*indexes = append(*indexes, nextIndex)
		*indexes = append(*indexes, nextIndex+2)
		*indexes = append(*indexes, nextIndex+3)
	}
	return nextIndex + 4
}

func (c *Chunk) freeModel() {
}

func (c *Chunk) buildRawModel(vertices []mgl32.Vec3, textures []mgl32.Vec2, normals []mgl32.Vec3, colors []mgl32.Vec3, indexes []uint32) {
	size := len(vertices)
	verticesArray := make([]float32, size*3)
	texturesArray := make([]float32, size*2)
	normalsArray := make([]float32, size*3)
	colorsArray := make([]float32, size*3)
	for i := 0; i < size; i++ {
		verticesArray[3*i] = vertices[i].X()
		verticesArray[3*i+1] = vertices[i].Y()
		verticesArray[3*i+2] = vertices[i].Z()
		texturesArray[2*i] = textures[i].X()
		texturesArray[2*i+1] = textures[i].Y()
		normalsArray[3*i] = normals[i].X()
		normalsArray[3*i+1] = normals[i].Y()
		normalsArray[3*i+2] = normals[i].Z()
		colorsArray[3*i] = colors[i].X()
		colorsArray[3*i+1] = colors[i].Y()
		colorsArray[3*i+2] = colors[i].Z()
	}
	if len(indexes) == 0 {
		c.Model.VertexCount = 0
	} else {
		c.Model = loader.LoadToVAO(verticesArray, texturesArray, indexes, normalsArray, colorsArray)
	}
}
