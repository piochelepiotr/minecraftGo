package worldcontent

import (
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
	blockType(x, y, z int, noises *noisesWithNeighbors) block.Block
	getStructures() []*structure
	getScale() int
	maxElevation() int
	minElevation() int
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

func noise3d(perlin *perlin.Perlin, x int, y int, z int, scale float64) float64 {
	c := perlin.Noise3D(float64(x)/scale, float64(y)/scale, float64(z)/scale)
	return (c + maxPerlin3D/2) / maxPerlin3D
}

type Generator struct {
	// perlin *perlin.Perlin
	noises *noisesLoader
	biomes []biome
	maxStructureSizeX int
	maxStructureSizeZ int
}

func makeBiomes(seed int64) []biome {
	biomes := make([]biome, 0)
	biomes = append(biomes, makePlainBiome(seed, 0))
	biomes = append(biomes, makeForestBiome(seed, 1))
	biomes = append(biomes, makeDesertBiome(seed, 2))
	return biomes
}

func newGenerator(worldConfig Config) *Generator {
	biomes:= makeBiomes(worldConfig.Seed)
	g := &Generator{
		biomes: biomes,
		noises: newNoisesLoader(worldConfig.Seed, biomes),
	}
	g.computeMaxStructureSizes()
	return g
}

// blockType returns block, without placing structures
func (g *Generator) blockType(x, y, z int, noises *noisesWithNeighbors) block.Block {
	biome := g.biomes[noises.getNoise(x, z).biome]
	// if structureBlock := g.getStructureBlock(biome, x, y, z, distanceFromBorder, noises); structureBlock != block.Air {
	// 	return structureBlock
	// }
	return biome.blockType(x, y, z, noises)
}

// generateChunkColumn allows you to create a chunk by passing the start point (the second chunk is at position ChunkSize-1)
func (g *Generator) generateChunkColumn(start geometry.Point2D) (chunks []*RawChunk) {
	n := WorldHeight/ChunkSize
	chunks = make([]*RawChunk, 0, n)
	noises := g.noises.getNoisesWithNeighbors(start[0], start[1])
	for y := 0; y < n; y ++ {
		chunk := RawChunk{}
		chunk.start = geometry.Point{start[0], y*ChunkSize, start[1]}
		chunk.blocks = make([]block.Block, ChunkSize*ChunkSize*ChunkSize)
		for x := 0; x < ChunkSize; x++ {
			for z := 0; z < ChunkSize; z++ {
				for y := 0; y < ChunkSize; y++ {
					chunk.blocks[index(x, y, z)] = g.blockType(chunk.start.X+x, chunk.start.Y+y, chunk.start.Z+z, noises)
				}
			}
		}
		chunks = append(chunks, &chunk)
	}
	g.addStructures(chunks, noises)
	return chunks
}

func (g *Generator) computeMaxStructureSizes() {
	var x, z int
	for _, b := range g.biomes {
		for _, s := range b.getStructures() {
			if s.x() > x {
				x = s.x()
			}
			if s.z() > z {
				z = s.z()
			}
		}
	}
	g.maxStructureSizeX = x
	g.maxStructureSizeZ = z
}

func (g *Generator) addStructures(chunks []*RawChunk, noises *noisesWithNeighbors) {
	startX := chunks[0].start.X
	startZ := chunks[0].start.Z
	setBlock := func(x, y, z int, b block.Block) {
		if x < startX || x >= startX + ChunkSize || z < startZ || z >= startZ + ChunkSize {
			return
		}
		x = x - startX
		z = z - startZ
		// works only because y > 0
		chunkY := y / ChunkSize
		y = y % ChunkSize
		chunks[chunkY].blocks[index(x, y, z)] = b
	}
	for x := startX - g.maxStructureSizeX; x < startX + ChunkSize + g.maxStructureSizeX; x++ {
		for z := startZ - g.maxStructureSizeZ; z < startZ + ChunkSize + g.maxStructureSizeZ; z++ {
			n := noises.getNoise(x, z)
			if n.structure == -1 {
				continue
			}
			s := g.biomes[n.biome].getStructures()[n.structure]
			xo := x + s.originX
			zo := z + s.originZ
			yo := noises.getNoise(xo, zo).elevation + 1
			for xs := 0; xs < s.x(); xs++ {
				for ys := 0; ys < s.y(); ys++ {
					for zs := 0; zs < s.z(); zs++ {
						setBlock(xs+x, ys+yo, zs+z, s.blocks[xs][ys][zs])
					}
				}
			}
		}
	}
}
