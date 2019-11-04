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
	TransparentModel  models.RawModel
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

type constructionChunk struct {
	vertices []mgl32.Vec3
	textures []mgl32.Vec2
	normals []mgl32.Vec3
	colors []mgl32.Vec3
	indices []uint32
	nextIndex uint32
}

func newConstructionChunk() *constructionChunk {
	return &constructionChunk{
		vertices: make([]mgl32.Vec3, 0),
		textures: make([]mgl32.Vec2, 0),
		normals: make([]mgl32.Vec3, 0),
		colors: make([]mgl32.Vec3, 0),
		indices: make([]uint32, 0),
	}
}

func (c *Chunk) buildFaces() {
	chunk := newConstructionChunk()
	transparentChunk := newConstructionChunk()
	for x := 0; x < ChunkSize; x++ {
		for y := 0; y < ChunkSize; y++ {
			for z := 0; z < ChunkSize; z++ {
				b := c.GetBlock(x, y, z)
				transparent := b.IsTransparent()
				//add face if the block isn't air and the block next to it is air
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
				if b != Air && (y == ChunkSize-1 || c.GetBlock(x, y+1, z).IsTransparent() || transparent) {
					n := mgl32.Vec3{0, 1, 0}
					p1 := mgl32.Vec3{xF, yF + 1, zF}
					p2 := mgl32.Vec3{xF + 1, yF + 1, zF}
					p3 := mgl32.Vec3{xF + 1, yF + 1, zF + 1}
					p4 := mgl32.Vec3{xF, yF + 1, zF + 1}
					if transparent {
						transparentChunk.addFace(p1, p2, p3, p4, n, b.GetSide(Top), true)
					} else {
						chunk.addFace(p1, p2, p3, p4, n, b.GetSide(Top), true)
					}
				}
				//bottom
				if b != Air && (y == 0 || c.GetBlock(x, y-1, z).IsTransparent() || transparent) {
					n := mgl32.Vec3{0, -1, 0}
					p1 := mgl32.Vec3{xF, yF, zF}
					p2 := mgl32.Vec3{xF + 1, yF, zF}
					p3 := mgl32.Vec3{xF + 1, yF, zF + 1}
					p4 := mgl32.Vec3{xF, yF, zF + 1}
					if transparent {
						transparentChunk.addFace(p1, p2, p3, p4, n, b.GetSide(Bottom), false)
					} else {
						chunk.addFace(p1, p2, p3, p4, n, b.GetSide(Bottom), false)
					}
				}
				//right
				if b != Air && (x == ChunkSize - 1 || c.GetBlock(x+1, y, z).IsTransparent() || transparent) {
					n := mgl32.Vec3{1, 0, 0}
					p1 := mgl32.Vec3{xF + 1, yF + 1, zF + 1}
					p2 := mgl32.Vec3{xF + 1, yF + 1, zF}
					p3 := mgl32.Vec3{xF + 1, yF, zF}
					p4 := mgl32.Vec3{xF + 1, yF, zF + 1}
					if transparent {
						transparentChunk.addFace(p1, p2, p3, p4, n, b.GetSide(Side), true)
					} else {
						chunk.addFace(p1, p2, p3, p4, n, b.GetSide(Side), true)
					}
				}
				//left
				if b != Air && (x == 0 || c.GetBlock(x-1, y, z).IsTransparent() || transparent) {
					n := mgl32.Vec3{-1, 0, 0}
					p1 := mgl32.Vec3{xF, yF + 1, zF}
					p2 := mgl32.Vec3{xF, yF + 1, zF + 1}
					p3 := mgl32.Vec3{xF, yF, zF + 1}
					p4 := mgl32.Vec3{xF, yF, zF}
					if transparent {
						transparentChunk.addFace(p1, p2, p3, p4, n, b.GetSide(Side), true)
					} else {
						chunk.addFace(p1, p2, p3, p4, n, b.GetSide(Side), true)
					}
				}
				//front
				if b != Air && (z == ChunkSize - 1 || c.GetBlock(x, y, z+1).IsTransparent() || transparent) {
					n := mgl32.Vec3{0, 0, 1}
					p1 := mgl32.Vec3{xF, yF + 1, zF + 1}
					p2 := mgl32.Vec3{xF + 1, yF + 1, zF + 1}
					p3 := mgl32.Vec3{xF + 1, yF, zF + 1}
					p4 := mgl32.Vec3{xF, yF, zF + 1}
					if transparent {
						transparentChunk.addFace(p1, p2, p3, p4, n, b.GetSide(Side), true)
					} else {
						chunk.addFace(p1, p2, p3, p4, n, b.GetSide(Side), true)
					}
				}
				//back
				if b != Air && (z == 0 || c.GetBlock(x, y, z-1).IsTransparent() || transparent) {
					n := mgl32.Vec3{0, 0, -1}
					p1 := mgl32.Vec3{xF + 1, yF + 1, zF}
					p2 := mgl32.Vec3{xF, yF + 1, zF}
					p3 := mgl32.Vec3{xF, yF, zF}
					p4 := mgl32.Vec3{xF + 1, yF, zF}
					if transparent {
						transparentChunk.addFace(p1, p2, p3, p4, n, b.GetSide(Side), true)
					} else {
						chunk.addFace(p1, p2, p3, p4, n, b.GetSide(Side), true)
					}
				}
			}
		}
	}
	if len(chunk.indices) > 0 {
		c.Model = loader.LoadToVAO(flatten3D(chunk.vertices), flatten2D(chunk.textures), chunk.indices, flatten3D(chunk.normals), flatten3D(chunk.colors))
	}
	if len(transparentChunk.indices) > 0 {
		c.TransparentModel = loader.LoadToVAO(flatten3D(transparentChunk.vertices), flatten2D(transparentChunk.textures), transparentChunk.indices, flatten3D(transparentChunk.normals), flatten3D(transparentChunk.colors))
	}
}

func (c *constructionChunk) addFace(p1, p2, p3, p4, n mgl32.Vec3, b Block, inverseRotation bool) {
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
	c.vertices = append(c.vertices, p1)
	c.vertices = append(c.vertices, p2)
	c.vertices = append(c.vertices, p3)
	c.vertices = append(c.vertices, p4)
	c.normals = append(c.normals, n)
	c.normals = append(c.normals, n)
	c.normals = append(c.normals, n)
	c.normals = append(c.normals, n)

	color := b.GetColor()
	c.colors = append(c.colors, color)
	c.colors = append(c.colors, color)
	c.colors = append(c.colors, color)
	c.colors = append(c.colors, color)
	c.textures = append(c.textures, t1)
	c.textures = append(c.textures, t2)
	c.textures = append(c.textures, t3)
	c.textures = append(c.textures, t4)
	if inverseRotation {
		c.indices = append(c.indices, c.nextIndex)
		c.indices = append(c.indices, c.nextIndex+2)
		c.indices = append(c.indices, c.nextIndex+1)
		c.indices = append(c.indices, c.nextIndex)
		c.indices = append(c.indices, c.nextIndex+3)
		c.indices = append(c.indices, c.nextIndex+2)
	} else {
		c.indices = append(c.indices, c.nextIndex)
		c.indices = append(c.indices, c.nextIndex+1)
		c.indices = append(c.indices, c.nextIndex+2)

		c.indices = append(c.indices, c.nextIndex)
		c.indices = append(c.indices, c.nextIndex+2)
		c.indices = append(c.indices, c.nextIndex+3)
	}
	c.nextIndex += 4
}

func (c *Chunk) freeModel() {
}

func flatten2D(array2D []mgl32.Vec2) []float32 {
	array := make([]float32, 0, len(array2D)*2)
	for _, p := range array2D {
		array = append(array, p.X())
		array = append(array, p.Y())
	}
	return array
}

func flatten3D(array3D []mgl32.Vec3) []float32 {
	array := make([]float32, 0, len(array3D)*3)
	for _, p := range array3D {
		array = append(array, p.X())
		array = append(array, p.Y())
		array = append(array, p.Z())
	}
	return array
}
