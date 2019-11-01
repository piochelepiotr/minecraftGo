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
	perlinScale float64 = 100
	// WorldHeight is the height of the world in blocks
	WorldHeight int = ChunkSize * 3
	dirtLayerThikness int = 5
	treeProbability float64 = 0.004
)

type Biome interface {
	blockType(x, y, z int) Block
}

type ForestBiome struct {
	structures []*Structure
	perlin *perlin.Perlin
}

func makeForestBiome() *ForestBiome {
	structures := make([]*Structure, 0)
	structures = append(structures, makeTree())
	return &ForestBiome{
		structures: structures,
		perlin:       perlin.NewPerlin(alpha, beta, perlinN, 233),
	}
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
	return biomes
}

func NewGenerator() *Generator {
	return &Generator{
		perlin:       perlin.NewPerlin(alpha, beta, perlinN, 233),
		biomes: makeBiomes(),
	}
}

func (f *ForestBiome) worldHeight(x, z int) int {
	half := WorldHeight/2
	return half + int(float64(half)*f.perlin.Noise2D(float64(x)/perlinScale, float64(z)/perlinScale))
}

func (g *Generator) getBiome(x, z int) Biome {
	return g.biomes[0]
}

func (g *Generator) BlockType(x, y, z int) Block {
	return g.getBiome(x, z).blockType(x, y, z)
}

func (f *ForestBiome) blockType(x, y, z int) Block {
	height := f.worldHeight(x, z)
	for _, s := range f.structures {
		xn := s.X()
		yn := s.Y()
		zn := s.Z()
		xo := int(math.Floor(float64(x)/float64(xn)))*xn
		zo := int(math.Floor(float64(z)/float64(zn)))*zn
		yo := f.worldHeight(xo, zo) + 1
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
