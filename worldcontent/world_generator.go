package worldcontent

import (
	"github.com/aquilax/go-perlin"
	"github.com/piochelepiotr/minecraftGo/geometry"
	"github.com/piochelepiotr/minecraftGo/world/block"
	"math"
)

const (
	alpha float64 = 1
	beta float64 = 2
	perlinN int = 3
	treeProbability float64 = 0.04
	biomeScale float64 = 200
)


type biome interface {
	blockType(x, y, z int) block.Block
	getStructures() []*structure
	worldHeight(x, z int) int
}

type structure struct {
	blocks [][][]block.Block
	p float64
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
		blocks[ix] = make([][]block.Block , y)
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

func (g *Generator) random(x int, z int, p float64) bool {
	n := int(g.seed)*x + z*int(g.seed) + int(g.seed) + x*z*int(g.seed)
	return n % int(1/p) == 0
}

type Generator struct {
	seed int64
	perlin *perlin.Perlin
	biomes []biome
}

func makeBiomes() []biome {
	biomes := make([]biome, 0)
	biomes = append(biomes, makeForestBiome())
	biomes = append(biomes, makeDesertBiome())
	return biomes
}

func newGenerator(worldConfig Config) *Generator {
	g := &Generator{
		seed: worldConfig.Seed,
		perlin:       perlin.NewPerlin(alpha, beta, perlinN, worldConfig.Seed),
		biomes: makeBiomes(),
	}
	return g
}

// returns a number between 0 and 1 generated using perlin noise
func perlinCoef(p *perlin.Perlin, x, z int, scale float64) float64 {
	// make sure we don't touch 1
	c := 0.5 + 0.5*p.Noise2D(float64(x)/scale, float64(z)/scale)
	if c >= 1 {
		return 0.999
	}
	if c < 0 {
		return 0
	}
	return c
}

func (g *Generator) getBiome(x, z int) biome {
	r := perlinCoef(g.perlin, x, z, biomeScale)
	incr := float64(1) / float64(len(g.biomes))
	n := int(math.Floor(r / incr))
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
		xo := int(math.Floor(float64(x)/float64(xn)))*xn
		zo := int(math.Floor(float64(z)/float64(zn)))*zn
		yo := b.worldHeight(xo+s.originX, zo+s.originZ) + 1
		if y < yo || y >= yo + yn {
			continue
		}
		// fmt.Printf("checking if a tree should be generated xo:%d, zo:%d, x:%d, z:%d\n", xo, zo, x, z)
		if !g.random(xo, zo, s.p) {
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
func (g *Generator) generateChunk(start geometry.Point, chunkSize int) (chunk RawChunk) {
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
