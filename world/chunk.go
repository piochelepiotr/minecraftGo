package world
import (
	"github.com/piochelepiotr/minecraftGo/models"
	"github.com/piochelepiotr/minecraftGo/renderEngine"
	"github.com/piochelepiotr/minecraftGo/textures"
	"github.com/go-gl/mathgl/mgl32"
)
type Chunk struct {
    Model models.RawModel
    blocks []Block
}
const NumberRowsTextures int = 2
const ChunkSize int = 16
const ChunkSize2 int = ChunkSize * ChunkSize
const ChunkSize3 int = ChunkSize2 * ChunkSize

func CreateChunk(startX int, startY int, startZ int, modelTexture textures.ModelTexture) Chunk{
    var chunk Chunk
    chunk.blocks = make([]Block, ChunkSize3)
    for x := 0; x < ChunkSize; x++ {
        for y := 0; y < ChunkSize; y++ {
            for z := 0; z < ChunkSize; z++ {
                chunk.SetBlock(x, y, z, Dirt)
            }
        }
    }
    chunk.buildFaces()
    return chunk
}

func (c *Chunk) SetBlock(x, y, z int, b Block) {
    c.blocks[x*ChunkSize2 + y*ChunkSize + z] = b
}

func (c *Chunk) GetBlock(x, y, z int) Block {
    return c.blocks[x*ChunkSize2 + y*ChunkSize + z]
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
                x_f := float32(x)
                y_f := float32(y)
                z_f := float32(z)
                if !(b == a || (y + 1 < ChunkSize && c.GetBlock(x, y+1, z) != a)) {
                    n := mgl32.Vec3{0,1,0}
                    p1 := mgl32.Vec3{x_f  , y_f+1, z_f  }
                    p2 := mgl32.Vec3{x_f+1, y_f+1, z_f  }
                    p3 := mgl32.Vec3{x_f+1, y_f+1, z_f+1}
                    p4 := mgl32.Vec3{x_f  , y_f+1, z_f+1}
                    c.addFace(&vertices, &textures, &normals, &indexes, p1, p2, p3, p4, n, b)
                }
                //bottom
                if !(b == a || (y - 1 > 0 && c.GetBlock(x, y-1, z) != a)) {
                    n := mgl32.Vec3{0,-1,0}
                    p1 := mgl32.Vec3{x_f  , y_f, (z_f  )}
                    p2 := mgl32.Vec3{x_f+1, y_f, (z_f  )}
                    p3 := mgl32.Vec3{x_f+1, y_f, (z_f+1)}
                    p4 := mgl32.Vec3{x_f  , y_f, (z_f+1)}
                    c.addFace(&vertices, &textures, &normals, &indexes, p1, p2, p3, p4, n, b)
                }
                //right
                if !(b == a || (x + 1 < ChunkSize && c.GetBlock(x+1, y, z) != a)) {
                    n := mgl32.Vec3{1,0,0}
                    p1 := mgl32.Vec3{(x_f+1), (y_f+1), (z_f+1)}
                    p2 := mgl32.Vec3{(x_f+1), (y_f+1), (z_f  )}
                    p3 := mgl32.Vec3{(x_f+1), (y_f  ), (z_f  )}
                    p4 := mgl32.Vec3{(x_f+1), (y_f  ), (z_f+1)}
                    c.addFace(&vertices, &textures, &normals, &indexes, p1, p2, p3, p4, n, b)
                }
                //left
                if !(b == a || (x - 1 > 0 && c.GetBlock(x-1, y, z) != a)) {
                    n := mgl32.Vec3{-1,0,0}
                    p1 := mgl32.Vec3{(x_f), (y_f+1), (z_f  )}
                    p2 := mgl32.Vec3{(x_f), (y_f+1), (z_f+1)}
                    p3 := mgl32.Vec3{(x_f), (y_f  ), (z_f+1)}
                    p4 := mgl32.Vec3{(x_f), (y_f  ), (z_f  )}
                    c.addFace(&vertices, &textures, &normals, &indexes, p1, p2, p3, p4, n, b)
                }
                //front
                if !(b == a || (z + 1 < ChunkSize && c.GetBlock(x, y, z+1) != a)) {
                    n := mgl32.Vec3{0,0,1}
                    p1 := mgl32.Vec3{(x_f  ), (y_f+1), (z_f+1)}
                    p2 := mgl32.Vec3{(x_f+1), (y_f+1), (z_f+1)}
                    p3 := mgl32.Vec3{(x_f+1), (y_f  ), (z_f+1)}
                    p4 := mgl32.Vec3{(x_f  ), (y_f  ), (z_f+1)}
                    c.addFace(&vertices, &textures, &normals, &indexes, p1, p2, p3, p4, n, b)
                }
                //back
                if !(b == a || (z - 1 > 0 && c.GetBlock(x, y, z-1) != a)) {
                    n := mgl32.Vec3{0,0,-1}
                    p1 := mgl32.Vec3{(x_f+1), (y_f+1), (z_f)}
                    p2 := mgl32.Vec3{(x_f  ), (y_f+1), (z_f)}
                    p3 := mgl32.Vec3{(x_f  ), (y_f  ), (z_f)}
                    p4 := mgl32.Vec3{(x_f+1), (y_f  ), (z_f)}
                    c.addFace(&vertices, &textures, &normals, &indexes, p1, p2, p3, p4, n, b)
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
    t1 := mgl32.Vec2{0,0}
    t2 := mgl32.Vec2{1,0}
    t3 := mgl32.Vec2{1,1}
    t4 := mgl32.Vec2{0,1}
    t1 = t1.Mul(1.0/float32(NumberRowsTextures)).Add(offsetTexture)
    t2 = t2.Mul(1.0/float32(NumberRowsTextures)).Add(offsetTexture)
    t3 = t3.Mul(1.0/float32(NumberRowsTextures)).Add(offsetTexture)
    t4 = t4.Mul(1.0/float32(NumberRowsTextures)).Add(offsetTexture)
    *vertices = append(*vertices, p1)
    *vertices = append(*vertices, p2)
    *vertices = append(*vertices, p3)
    *vertices = append(*vertices, p4)
    *normals =  append(*normals, n)
    *normals =  append(*normals, n)
    *normals =  append(*normals, n)
    *normals =  append(*normals, n)
    *textures = append(*textures, t1)
    *textures = append(*textures, t2)
    *textures = append(*textures, t3)
    *textures = append(*textures, t4)
    start_index := uint32(0)
    if len(*indexes) > 0 {
        start_index = (*indexes)[len(*indexes) - 1] + 1
    }
    *indexes = append(*indexes, start_index)
    *indexes = append(*indexes, start_index+1)
    *indexes = append(*indexes, start_index+2)

    *indexes = append(*indexes, start_index)
    *indexes = append(*indexes, start_index+2)
    *indexes = append(*indexes, start_index+3)
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
    c.Model = renderEngine.LoadToVAO(verticesArray, texturesArray, indexes, normalsArray)
}

