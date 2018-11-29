package world

import (
	"fmt"

	"github.com/aquilax/go-perlin"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/piochelepiotr/minecraftGo/loader"
	"github.com/piochelepiotr/minecraftGo/models"
	"github.com/piochelepiotr/minecraftGo/textures"
)

// Chunk is set cube of blocks
type Chunk struct {
	Model  models.RawModel
	blocks []Block
	perlin *perlin.Perlin
	Start  Point
}

// NumberRowsTextures is the number number of rows on the texture image
const NumberRowsTextures int = 2

// ChunkSize is the size of a chunk in blocks
const ChunkSize int = 16

// ChunkSize2 is the area of a chunk in blocks
const ChunkSize2 int = ChunkSize * ChunkSize

// ChunkSize3 is the volume of a chunk in blocks
const ChunkSize3 int = ChunkSize2 * ChunkSize
const alpha float64 = 2
const beta float64 = 2
const perlinN int = 3
const perlinScale float64 = 20

// WorldHeight is the height of the world in blocks
const WorldHeight int = ChunkSize * 2

// CreateChunk allows you to create a chunk by passing the start point (the second chunk is at position ChunkSize-1)
func CreateChunk(startX int, startY int, startZ int, modelTexture textures.ModelTexture, perlin *perlin.Perlin) Chunk {
	var chunk Chunk
	chunk.Start = Point{
		X: startX,
		Y: startY,
		Z: startZ,
	}
	chunk.blocks = make([]Block, ChunkSize3)
	chunk.perlin = perlin
	halfWorldHeight := WorldHeight / 2
	for x := 0; x < ChunkSize; x++ {
		for z := 0; z < ChunkSize; z++ {
			height := halfWorldHeight + int(float64(halfWorldHeight)*chunk.perlin.Noise2D(float64(startX+x)/perlinScale, float64(startZ+z)/perlinScale)) - startY
			positiveHeight := height
			if positiveHeight < 0 {
				positiveHeight = 0
			}
			for y := 0; y < positiveHeight && y < ChunkSize; y++ {
				chunk.setBlockNoUpdate(x, y, z, Dirt)
			}
			for y := positiveHeight; y < ChunkSize; y++ {
				chunk.setBlockNoUpdate(x, y, z, Air)
			}
			if 0 <= height && height < ChunkSize {
				chunk.setBlockNoUpdate(x, height, z, Grass)
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
	indexes := make([]uint32, 0)
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
					c.addFace(&vertices, &textures, &normals, &indexes, p1, p2, p3, p4, n, BlockSides[b]["top"])
				}
				//bottom
				if !(b == a || (y-1 > 0 && c.GetBlock(x, y-1, z) != a)) {
					n := mgl32.Vec3{0, -1, 0}
					p1 := mgl32.Vec3{xF, yF, (zF)}
					p2 := mgl32.Vec3{xF + 1, yF, (zF)}
					p3 := mgl32.Vec3{xF + 1, yF, (zF + 1)}
					p4 := mgl32.Vec3{xF, yF, (zF + 1)}
					c.addFace(&vertices, &textures, &normals, &indexes, p1, p2, p3, p4, n, BlockSides[b]["bottom"])
				}
				//right
				if !(b == a || (x+1 < ChunkSize && c.GetBlock(x+1, y, z) != a)) {
					n := mgl32.Vec3{1, 0, 0}
					p1 := mgl32.Vec3{(xF + 1), (yF + 1), (zF + 1)}
					p2 := mgl32.Vec3{(xF + 1), (yF + 1), (zF)}
					p3 := mgl32.Vec3{(xF + 1), (yF), (zF)}
					p4 := mgl32.Vec3{(xF + 1), (yF), (zF + 1)}
					c.addFace(&vertices, &textures, &normals, &indexes, p1, p2, p3, p4, n, BlockSides[b]["side"])
				}
				//left
				if !(b == a || (x-1 > 0 && c.GetBlock(x-1, y, z) != a)) {
					n := mgl32.Vec3{-1, 0, 0}
					p1 := mgl32.Vec3{(xF), (yF + 1), (zF)}
					p2 := mgl32.Vec3{(xF), (yF + 1), (zF + 1)}
					p3 := mgl32.Vec3{(xF), (yF), (zF + 1)}
					p4 := mgl32.Vec3{(xF), (yF), (zF)}
					c.addFace(&vertices, &textures, &normals, &indexes, p1, p2, p3, p4, n, BlockSides[b]["side"])
				}
				//front
				if !(b == a || (z+1 < ChunkSize && c.GetBlock(x, y, z+1) != a)) {
					n := mgl32.Vec3{0, 0, 1}
					p1 := mgl32.Vec3{(xF), (yF + 1), (zF + 1)}
					p2 := mgl32.Vec3{(xF + 1), (yF + 1), (zF + 1)}
					p3 := mgl32.Vec3{(xF + 1), (yF), (zF + 1)}
					p4 := mgl32.Vec3{(xF), (yF), (zF + 1)}
					c.addFace(&vertices, &textures, &normals, &indexes, p1, p2, p3, p4, n, BlockSides[b]["side"])
				}
				//back
				if !(b == a || (z-1 > 0 && c.GetBlock(x, y, z-1) != a)) {
					n := mgl32.Vec3{0, 0, -1}
					p1 := mgl32.Vec3{(xF + 1), (yF + 1), (zF)}
					p2 := mgl32.Vec3{(xF), (yF + 1), (zF)}
					p3 := mgl32.Vec3{(xF), (yF), (zF)}
					p4 := mgl32.Vec3{(xF + 1), (yF), (zF)}
					c.addFace(&vertices, &textures, &normals, &indexes, p1, p2, p3, p4, n, BlockSides[b]["side"])
				}
			}
		}
	}
	c.buildRawModel(vertices, textures, normals, indexes)
}

func (c *Chunk) addFace(vertices *[]mgl32.Vec3, textures *[]mgl32.Vec2, normals *[]mgl32.Vec3, indexes *[]uint32, p1, p2, p3, p4, n mgl32.Vec3, b Block) {
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
	*textures = append(*textures, t1)
	*textures = append(*textures, t2)
	*textures = append(*textures, t3)
	*textures = append(*textures, t4)
	startIndex := uint32(0)
	if len(*indexes) > 0 {
		startIndex = (*indexes)[len(*indexes)-1] + 1
	}
	*indexes = append(*indexes, startIndex)
	*indexes = append(*indexes, startIndex+1)
	*indexes = append(*indexes, startIndex+2)

	*indexes = append(*indexes, startIndex)
	*indexes = append(*indexes, startIndex+2)
	*indexes = append(*indexes, startIndex+3)
}

func (c *Chunk) freeModel() {
}

func (c *Chunk) buildRawModel(vertices []mgl32.Vec3, textures []mgl32.Vec2, normals []mgl32.Vec3, indexes []uint32) {
	size := len(vertices)
	verticesArray := make([]float32, size*3)
	texturesArray := make([]float32, size*2)
	normalsArray := make([]float32, size*3)
	for i := 0; i < size; i++ {
		verticesArray[3*i] = vertices[i].X()
		verticesArray[3*i+1] = vertices[i].Y()
		verticesArray[3*i+2] = vertices[i].Z()
		texturesArray[2*i] = textures[i].X()
		texturesArray[2*i+1] = textures[i].Y()
		normalsArray[3*i] = normals[i].X()
		normalsArray[3*i+1] = normals[i].Y()
		normalsArray[3*i+2] = normals[i].Z()
	}
	if len(indexes) == 0 {
		c.Model.VertexCount = 0
	} else {
		fmt.Println("building!!!")
		fmt.Println(c.Model.VaoID)
		c.Model = loader.LoadToVAO(verticesArray, texturesArray, indexes, normalsArray)
		fmt.Println(c.Model.VaoID)
	}
}
