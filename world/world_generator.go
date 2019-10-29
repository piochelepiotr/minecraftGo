package world

import "github.com/aquilax/go-perlin"

const (
	alpha float64 = 1
	beta float64 = 2
	perlinN int = 3
	perlinScale float64 = 100
	// WorldHeight is the height of the world in blocks
	WorldHeight int = ChunkSize * 3
	dirtLayerThikness int = 5
)

type Generator struct {
	perlin *perlin.Perlin
}

func NewGenerator() *Generator {
	return &Generator{
		perlin:       perlin.NewPerlin(alpha, beta, perlinN, 233),
	}
}

func (g *Generator) WorldHeight(x, z int) int {
	half := WorldHeight/2
	return half + int(float64(half)*g.perlin.Noise2D(float64(x)/perlinScale, float64(z)/perlinScale))
}

func (g *Generator) BlockType(x, y, z int) Block {
	height := g.WorldHeight(x, z)
	if y == height {
		return Grass
	}
	if y < height && y > height - dirtLayerThikness {
		return Dirt
	}
	if y <= height - dirtLayerThikness {
		return Stone
	}
	return Air
}
