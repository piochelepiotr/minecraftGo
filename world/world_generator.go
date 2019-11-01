package world

import (
	"github.com/aquilax/go-perlin"
	"math"
	"math/rand"
)

const (
	alpha float64 = 1
	beta float64 = 2
	perlinN int = 3
	// WorldHeight is the height of the world in blocks
	WorldHeight = ChunkSize * 3
	treeProbability float64 = 0.004
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
	s := makeStructure(3, 5, 3)
	s.blocks[1][0][1] = Tree
	s.blocks[1][1][1] = Tree
	s.blocks[1][2][1] = Tree
	s.blocks[1][3][1] = Tree
	s.blocks[1][4][1] = Leaves
	s.blocks[1][3][0] = Leaves
	s.blocks[1][3][2] = Leaves
	s.blocks[0][3][1] = Leaves
	s.blocks[2][3][1] = Leaves
	s.p = treeProbability
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

func NewGenerator() *Generator {
	return &Generator{
		perlin:       perlin.NewPerlin(alpha, beta, perlinN, 233),
		biomes: makeBiomes(),
	}
}

// returns a number between 0 and 1 generated using perlin noise
func perlinCoef(p *perlin.Perlin, x, z int, scale float64) float64 {
	return 0.5 + 0.5*p.Noise2D(float64(x)/scale, float64(z)/scale)
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
		yo := b.worldHeight(xo, zo) + 1
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
