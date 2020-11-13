package world

import (
	"github.com/aquilax/go-perlin"
	"github.com/piochelepiotr/minecraftGo/geometry"
	"math"
	"math/rand"
)

const (
	alpha float64 = 1
	beta float64 = 2
	perlinN int = 3
	// WorldHeight is the height of the world in blocks
	WorldHeight = ChunkSize * 3
	treeProbability float64 = 0.04
	biomeScale float64 = 200
)

type Biome interface {
	blockType(x, y, z int) Block
	getStructures() []*Structure
	worldHeight(x, z int) int
}

type Structure struct {
	blocks [][][]Block
	p float64
	originX int
	originZ int
}

func (s *Structure) X() int {
	return len(s.blocks)
}

func (s *Structure) Y() int {
	return len(s.blocks[0])
}

func (s *Structure) Z() int {
	return len(s.blocks[0][0])
}

func makeStructure(x, y, z int) *Structure {
	blocks := make([][][]Block, x)
	for ix := 0; ix < x; ix++ {
		blocks[ix] = make([][]Block , y)
		for iy := 0; iy < y; iy++ {
			blocks[ix][iy] = make([]Block, z)
			for iz := 0; iz < z; iz++ {
				blocks[ix][iy][iz] = Air
			}
		}
	}
	return &Structure{
		blocks: blocks,
	}
}

func makeTree() *Structure {
	s := makeStructure(5, 7, 5)
	for x := 0; x < 5; x++ {
		for y := 3; y < 5; y++ {
			for z := 0; z < 5; z++ {
				s.blocks[x][y][z] = Leaves
			}
		}
	}
	for x := 1; x < 4; x++ {
		for z := 1; z < 4; z++ {
			s.blocks[x][5][z] = Leaves
		}
	}
	for y := 0; y < 6; y++ {
		s.blocks[2][y][2] = Tree
	}
	for i := 1; i < 4; i++ {
		s.blocks[i][6][2] = Leaves
		s.blocks[2][6][i] = Leaves
	}
	s.p = treeProbability
	s.originX = 2
	s.originZ = 2
	return s
}

func random(x int, z int, p float64) bool {
	rand.Seed(int64((x + z)*z))
	return rand.Intn(int(1/p)) == 0
}

type Generator struct {
	perlin *perlin.Perlin
	biomes []Biome
}

func makeBiomes() []Biome {
	biomes := make([]Biome, 0)
	biomes = append(biomes, makeForestBiome())
	biomes = append(biomes, makeDesertBiome())
	return biomes
}

func NewGenerator(worldConfig Config) *Generator {
	return &Generator{
		perlin:       perlin.NewPerlin(alpha, beta, perlinN, worldConfig.Seed),
		biomes: makeBiomes(),
	}
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

func (g *Generator) getBiome(x, z int) Biome {
	r := perlinCoef(g.perlin, x, z, biomeScale)
	incr := float64(1) / float64(len(g.biomes))
	n := int(math.Floor(r / incr))
	return g.biomes[n]
}

func (g *Generator) BlockType(x, y, z int) Block {
	biome := g.getBiome(x, z)
	if structureBlock := getStructureBlock(biome, x, y, z); structureBlock != Air {
		return structureBlock
	}
	return biome.blockType(x, y, z)
}

func getStructureBlock(b Biome, x, y, z int) Block {
	for _, s := range b.getStructures() {
		xn := s.X()
		yn := s.Y()
		zn := s.Z()
		xo := int(math.Floor(float64(x)/float64(xn)))*xn
		zo := int(math.Floor(float64(z)/float64(zn)))*zn
		yo := b.worldHeight(xo+s.originX, zo+s.originZ) + 1
		if y < yo || y >= yo + yn {
			continue
		}
		// fmt.Printf("checking if a tree should be generated xo:%d, zo:%d, x:%d, z:%d\n", xo, zo, x, z)
		if !random(xo, zo, s.p) {
			continue
		}
		xi := x - xo
		yi := y - yo
		zi := z - zo
		return s.blocks[xi][yi][zi]
	}
	return Air
}

// GenerateChunk allows you to create a chunk by passing the start point (the second chunk is at position ChunkSize-1)
func (g *Generator) GenerateChunk(start geometry.Point) (chunk RawChunk) {
	chunk.size = ChunkSize
	chunk.Start = start
	chunk.blocks = make([]Block, ChunkSize*ChunkSize*ChunkSize)
	for x := 0; x < ChunkSize; x++ {
		for z := 0; z < ChunkSize; z++ {
			for y := 0; y < ChunkSize; y++ {
				chunk.setBlockNoUpdate(x, y, z, g.BlockType(start.X+x, start.Y+y, start.Z+z))
			}
		}
	}
	return chunk
}

// GetHeight returns height of the world in blocks at a x,z position
func (g *Generator) GetHeight(x, z int) int {
	for y := WorldHeight - 1; y >= 0; y-- {
		if g.BlockType(x, y, z) != Air {
			return y + 1
		}
	}
	return 0
}
