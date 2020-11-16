package worldcontent

import (
	"github.com/piochelepiotr/minecraftGo/random"
	"math"

	"github.com/aquilax/go-perlin"
	"github.com/piochelepiotr/minecraftGo/geometry"
	"github.com/piochelepiotr/minecraftGo/world/block"
)

const (
	treeProbability float64 = 0.5
	biomeScale      float64 = 200
)

var maxPerlin2D = math.Sqrt(2)
var maxPerlin3D = math.Sqrt(3)

type biome interface {
	blockType(x, y, z int) block.Block
	getStructures() []*structure
	worldHeight(x, z int) int
}

type structure struct {
	blocks  [][][]block.Block
	p       float64
	originX int
	originZ int
}

func (s *structure) x() int {
	return len(s.blocks)
}

func (s *structure) y() int {
	return len(s.blocks[0])
}

func (s *structure) z() int {
	return len(s.blocks[0][0])
}

func makeStructure(x, y, z int) *structure {
	blocks := make([][][]block.Block, x)
	for ix := 0; ix < x; ix++ {
		blocks[ix] = make([][]block.Block, y)
		for iy := 0; iy < y; iy++ {
			blocks[ix][iy] = make([]block.Block, z)
			for iz := 0; iz < z; iz++ {
				blocks[ix][iy][iz] = block.Air
			}
		}
	}
	return &structure{
		blocks: blocks,
	}
}

func makeTree() *structure {
	s := makeStructure(5, 7, 5)
	for x := 0; x < 5; x++ {
		for y := 3; y < 5; y++ {
			for z := 0; z < 5; z++ {
				s.blocks[x][y][z] = block.Leaves
			}
		}
	}
	for x := 1; x < 4; x++ {
		for z := 1; z < 4; z++ {
			s.blocks[x][5][z] = block.Leaves
		}
	}
	for y := 0; y < 6; y++ {
		s.blocks[2][y][2] = block.Tree
	}
	for i := 1; i < 4; i++ {
		s.blocks[i][6][2] = block.Leaves
		s.blocks[2][6][i] = block.Leaves
	}
	s.p = treeProbability
	s.originX = 2
	s.originZ = 2
	return s
}

func noise2d(p *perlin.Perlin, x int, y int, scale float64) float64 {
	c := p.Noise2D(float64(x)/scale, float64(y)/scale)
	c = (c + maxPerlin2D/2) / maxPerlin2D
	return c
}

func elevation(p *perlin.Perlin, x int, y int, scale float64, min int, max int) int {
	return min + int(float64(max-min)*noise2d(p, x, y, scale))
}

func noise3d(perlin *perlin.Perlin, x int, y int, z int, scale float64) float64 {
	c := perlin.Noise3D(float64(x)/scale, float64(y)/scale, float64(z)/scale)
	return (c + maxPerlin3D/2) / maxPerlin3D
}

func random2d(perlin *perlin.Perlin, x int, y int, p float64) bool {
	c := perlin.Noise2D(float64(x)+0.5, float64(y)+0.5)
	c = (c + maxPerlin2D/2) / maxPerlin2D
	return c <= p
}

func (g *Generator) random(x int, z int, p float64) bool {
	return g.noise.Noise2D(x, z) <= p
}

type Generator struct {
	seed   int64
	perlin *perlin.Perlin
	noise *random.Noise
	biomes []biome
}

func makeBiomes(seed int64) []biome {
	biomes := make([]biome, 0)
	biomes = append(biomes, makeForestBiome(seed))
	// biomes = append(biomes, makeDesertBiome())
	return biomes
}

func newGenerator(worldConfig Config) *Generator {
	g := &Generator{
		seed:   worldConfig.Seed,
		perlin: perlin.NewPerlin(2, 2, 1, worldConfig.Seed),
		noise: random.NewNoise(worldConfig.Seed),
		biomes: makeBiomes(worldConfig.Seed),
	}
	return g
}

func (g *Generator) getBiome(x, z int) biome {
	// change that
	n := elevation(g.perlin, x, z, biomeScale, 0, len(g.biomes))
	return g.biomes[n]
}

func (g *Generator) blockType(x, y, z int) block.Block {
	biome := g.getBiome(x, z)
	if structureBlock := g.getStructureBlock(biome, x, y, z); structureBlock != block.Air {
		return structureBlock
	}
	return biome.blockType(x, y, z)
}

func (g *Generator) getStructureBlock(b biome, x, y, z int) block.Block {
	for _, s := range b.getStructures() {
		xn := s.x()
		yn := s.y()
		zn := s.z()
		xo := int(math.Floor(float64(x)/float64(xn))) * xn
		zo := int(math.Floor(float64(z)/float64(zn))) * zn
		yo := b.worldHeight(xo+s.originX, zo+s.originZ) + 1
		if y < yo || y >= yo+yn {
			continue
		}
		// fmt.Printf("checking if a tree should be generated xo:%d, zo:%d, x:%d, z:%d\n", xo, zo, x, z)
		p := noise2d(g.perlin, xo, zo, 100)*s.p
		if !g.random(xo, zo, p) {
			continue
		}
		xi := x - xo
		yi := y - yo
		zi := z - zo
		return s.blocks[xi][yi][zi]
	}
	return block.Air
}

// generateChunk allows you to create a chunk by passing the start point (the second chunk is at position ChunkSize-1)
func (g *Generator) generateChunk(start geometry.Point) (chunk *RawChunk) {
	chunk = &RawChunk{}
	chunk.start = start
	chunk.blocks = make([]block.Block, ChunkSize*ChunkSize*ChunkSize)
	for x := 0; x < ChunkSize; x++ {
		for z := 0; z < ChunkSize; z++ {
			for y := 0; y < ChunkSize; y++ {
				chunk.blocks[index(x, y, z)] = g.blockType(start.X+x, start.Y+y, start.Z+z)
			}
		}
	}
	return chunk
}

// getHeight returns height of the world in blocks at a x,z position
func (g *Generator) getHeight(x, z, worldHeight int) int {
	for y := worldHeight - 1; y >= 0; y-- {
		if g.blockType(x, y, z) != block.Air {
			return y + 1
		}
	}
	return 0
}
