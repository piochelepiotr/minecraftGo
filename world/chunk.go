package world

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/piochelepiotr/minecraftGo/game_engine/loader"
	"github.com/piochelepiotr/minecraftGo/game_engine/models"
	"github.com/piochelepiotr/minecraftGo/geometry"
	"github.com/piochelepiotr/minecraftGo/world/block"
	"github.com/piochelepiotr/minecraftGo/worldcontent"
)

// numberRowsTextures is the number number of rows on the texture image
const numberRowsTextures int = 16

// Chunk is set cube of blocks
type Chunk struct {
	start geometry.Point
	model            *constructionChunk
	transparentModel *constructionChunk
	Model            models.RawModel
	TransparentModel models.RawModel
	world *worldcontent.InMemoryWorld
}

// NewChunk returns a new graphic chunk
func NewChunk(world *worldcontent.InMemoryWorld, start geometry.Point) *Chunk {
	chunk := Chunk{world: world, start: start}
	chunk.buildFaces()
	return &chunk
}

type constructionChunk struct {
	vertices  []mgl32.Vec3
	textures  []mgl32.Vec2
	normals   []mgl32.Vec3
	colors    []mgl32.Vec3
	indices   []uint32
	nextIndex uint32
}

func newConstructionChunk() *constructionChunk {
	return &constructionChunk{
		vertices: make([]mgl32.Vec3, 0),
		textures: make([]mgl32.Vec2, 0),
		normals:  make([]mgl32.Vec3, 0),
		colors:   make([]mgl32.Vec3, 0),
		indices:  make([]uint32, 0),
	}
}

// Load chunk to openGL
func (c *Chunk) Load() {
	if len(c.model.indices) > 0 {
		c.Model = loader.LoadToVAO(flatten3D(c.model.vertices), flatten2D(c.model.textures), c.model.indices, flatten3D(c.model.normals), flatten3D(c.model.colors))
	}
	if len(c.transparentModel.indices) > 0 {
		c.TransparentModel = loader.LoadToVAO(flatten3D(c.transparentModel.vertices), flatten2D(c.transparentModel.textures), c.transparentModel.indices, flatten3D(c.transparentModel.normals), flatten3D(c.transparentModel.colors))
	}
}

func (c *constructionChunk) buildBlock(x, y, z float32, b block.Block, up, bottom, right, left, front, back bool) {
	//add face if the block next to it is transparent
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
	if up {
		n := mgl32.Vec3{0, 1, 0}
		p1 := mgl32.Vec3{x, y + 1, z}
		p2 := mgl32.Vec3{x + 1, y + 1, z}
		p3 := mgl32.Vec3{x + 1, y + 1, z + 1}
		p4 := mgl32.Vec3{x, y + 1, z + 1}
		c.addFace(p1, p2, p3, p4, n, b.GetSide(block.Top), true)
	}
	//bottom
	if bottom {
		n := mgl32.Vec3{0, -1, 0}
		p1 := mgl32.Vec3{x, y, z}
		p2 := mgl32.Vec3{x + 1, y, z}
		p3 := mgl32.Vec3{x + 1, y, z + 1}
		p4 := mgl32.Vec3{x, y, z + 1}
		c.addFace(p1, p2, p3, p4, n, b.GetSide(block.Bottom), false)
	}
	//right
	if right {
		n := mgl32.Vec3{1, 0, 0}
		p1 := mgl32.Vec3{x + 1, y + 1, z + 1}
		p2 := mgl32.Vec3{x + 1, y + 1, z}
		p3 := mgl32.Vec3{x + 1, y, z}
		p4 := mgl32.Vec3{x + 1, y, z + 1}
		c.addFace(p1, p2, p3, p4, n, b.GetSide(block.Side), true)
	}
	//left
	if left {
		n := mgl32.Vec3{-1, 0, 0}
		p1 := mgl32.Vec3{x, y + 1, z}
		p2 := mgl32.Vec3{x, y + 1, z + 1}
		p3 := mgl32.Vec3{x, y, z + 1}
		p4 := mgl32.Vec3{x, y, z}
		c.addFace(p1, p2, p3, p4, n, b.GetSide(block.Side), true)
	}
	//front
	if front {
		n := mgl32.Vec3{0, 0, 1}
		p1 := mgl32.Vec3{x, y + 1, z + 1}
		p2 := mgl32.Vec3{x + 1, y + 1, z + 1}
		p3 := mgl32.Vec3{x + 1, y, z + 1}
		p4 := mgl32.Vec3{x, y, z + 1}
		c.addFace(p1, p2, p3, p4, n, b.GetSide(block.Side), true)
	}
	//back
	if back {
		n := mgl32.Vec3{0, 0, -1}
		p1 := mgl32.Vec3{x + 1, y + 1, z}
		p2 := mgl32.Vec3{x, y + 1, z}
		p3 := mgl32.Vec3{x, y, z}
		p4 := mgl32.Vec3{x + 1, y, z}
		c.addFace(p1, p2, p3, p4, n, b.GetSide(block.Side), true)
	}
}

func (c *Chunk) buildFaces() {
	c.model = newConstructionChunk()
	c.transparentModel = newConstructionChunk()
	for x := 0; x < c.world.ChunkSize(); x++ {
		for y := 0; y < c.world.ChunkSize(); y++ {
			for z := 0; z < c.world.ChunkSize(); z++ {
				b := c.world.GetBlock(c.start.X + x, c.start.Y + y, c.start.Z+ z)
				if b == block.Air {
					continue
				}
				up := c.world.GetBlock(c.start.X + x, c.start.Y + y+1, c.start.Z + z).IsTransparent()
				bottom := c.world.GetBlock(c.start.X + x, c.start.Y + y-1, c.start.Z + z).IsTransparent()
				right := c.world.GetBlock(c.start.X + x+1, c.start.Y + y, c.start.Z + z).IsTransparent()
				left := c.world.GetBlock(c.start.X + x-1, c.start.Y + y, c.start.Z + z).IsTransparent()
				front := c.world.GetBlock(c.start.X + x, c.start.Y + y, c.start.Z + z+1).IsTransparent()
				back := c.world.GetBlock(c.start.X + x, c.start.Y + y, c.start.Z + z-1).IsTransparent()
				cons := c.model
				if b.IsTransparent() {
					cons = c.transparentModel
				}
				cons.buildBlock(float32(x), float32(y), float32(z), b, up, bottom, right, left, front, back)
			}
		}
	}
}

func (c *constructionChunk) addFace(p1, p2, p3, p4, n mgl32.Vec3, b block.Block, inverseRotation bool) {
	if b == block.Grass {
		if n.ApproxEqual(mgl32.Vec3{0, -1, 0}) {
			b = block.Dirt
		} else if !n.ApproxEqual(mgl32.Vec3{0, 1, 0}) {
			b = block.Grass
		}
	}
	textureX := int(b) % numberRowsTextures
	textureY := int(b) / numberRowsTextures
	offsetTextureX := float32(textureX) / float32(numberRowsTextures)
	offsetTextureY := float32(textureY) / float32(numberRowsTextures)
	offsetTexture := mgl32.Vec2{offsetTextureX, offsetTextureY}
	t1 := mgl32.Vec2{0, 0}
	t2 := mgl32.Vec2{1, 0}
	t3 := mgl32.Vec2{1, 1}
	t4 := mgl32.Vec2{0, 1}
	t1 = t1.Mul(1.0 / float32(numberRowsTextures)).Add(offsetTexture)
	t2 = t2.Mul(1.0 / float32(numberRowsTextures)).Add(offsetTexture)
	t3 = t3.Mul(1.0 / float32(numberRowsTextures)).Add(offsetTexture)
	t4 = t4.Mul(1.0 / float32(numberRowsTextures)).Add(offsetTexture)
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
