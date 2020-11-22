package worldcontent

import (
	"github.com/piochelepiotr/minecraftGo/random"
	"math"

	"github.com/aquilax/go-perlin"
	"github.com/piochelepiotr/minecraftGo/geometry"
	"github.com/piochelepiotr/minecraftGo/world/block"
)

const (
	biomeScale      float64 = 300
	nBiomes = 3
)

var maxPerlin2D = math.Sqrt(2)
var maxPerlin3D = math.Sqrt(3)

type biome interface {
	blockType(x, y, z int, distanceFromBorder float64, noises *noisesWithNeighbors) block.Block
	getStructures() []*structure
	worldHeight(x, z int, distanceFromBorder float64, noises *noisesWithNeighbors) int
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

func elevation(noise float64, min int, max int, distanceToBorder float64) int {
	return min + int(float64(max-min)*noise*distanceToBorder)
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
	// perlin *perlin.Perlin
	noises *noisesLoader
	noise *random.Noise
	biomes []biome
}

func makeBiomes(seed int64) []biome {
	biomes := make([]biome, 0)
	biomes = append(biomes, makePlainBiome(seed, 0))
	biomes = append(biomes, makeForestBiome(seed, 1))
	biomes = append(biomes, makeDesertBiome(seed, 2))
	return biomes
}

func newGenerator(worldConfig Config) *Generator {
	g := &Generator{
		noise: random.NewNoise(worldConfig.Seed),
		biomes: makeBiomes(worldConfig.Seed),
		noises: newNoisesLoader(worldConfig.Seed),
	}
	return g
}

// if distance != 1, we are at the border
const borderSize = 0.4
func distanceFromBiomeBorder(p float64) float64 {
	lower := float64(int(p))
	upper := lower + 1
	toLower := (p - lower)*(1/borderSize)
	toUpper := (upper - p)*(1/borderSize)
	minDistance := math.Min(toLower, toUpper)
	return math.Min(minDistance, 1)
}

func (g *Generator) getBiome(x, z int, noises *noisesWithNeighbors) (biome biome, distanceFromBorder float64) {
	// change that
	p := noises.getNoise(x, z).biomeNoise * float64(len(g.biomes))
	return g.biomes[int(p)], distanceFromBiomeBorder(p)
}

func (g *Generator) blockType(x, y, z int, noises *noisesWithNeighbors) block.Block {
	biome, distanceFromBorder := g.getBiome(x, z, noises)
	if structureBlock := g.getStructureBlock(biome, x, y, z, distanceFromBorder, noises); structureBlock != block.Air {
		return structureBlock
	}
	return biome.blockType(x, y, z, distanceFromBorder, noises)
}

func (g *Generator) getStructureBlock(b biome, x, y, z int, distanceFromBorder float64, noises *noisesWithNeighbors) block.Block {
	for _, s := range b.getStructures() {
		xn := s.x()
		yn := s.y()
		zn := s.z()
		xo := int(math.Floor(float64(x)/float64(xn))) * xn
		zo := int(math.Floor(float64(z)/float64(zn))) * zn
		yo := b.worldHeight(xo+s.originX, zo+s.originZ, distanceFromBorder, noises) + 1
		if y < yo || y >= yo+yn {
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
func (g *Generator) generateChunk(start geometry.Point) (chunk *RawChunk) {
	chunk = &RawChunk{}
	chunk.start = start
	chunk.blocks = make([]block.Block, ChunkSize*ChunkSize*ChunkSize)
	noises := g.noises.getNoisesWithNeighbors(start.X, start.Z)
	for x := 0; x < ChunkSize; x++ {
		for z := 0; z < ChunkSize; z++ {
			for y := 0; y < ChunkSize; y++ {
				chunk.blocks[index(x, y, z)] = g.blockType(start.X+x, start.Y+y, start.Z+z, noises)
			}
		}
	}
	return chunk
}