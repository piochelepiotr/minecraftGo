package worldcontent

import (
	"github.com/aquilax/go-perlin"
	"github.com/piochelepiotr/minecraftGo/world/block"
)

const (
	desertMinElevation = 60
	// max is not reached, max - 1 is reached
	desertMaxElevation = 65
	desertElevationScale = 100
)

type DesertBiome struct {
	index int
	structures []*structure
	perlin *perlin.Perlin
}

func (d *DesertBiome) getStructures() []*structure {
	return d.structures
}

func makeDesertBiome(globalSeed int64, index int) *DesertBiome {
	seed := 5*globalSeed
	structures := make([]*structure, 0)
	structures = append(structures, makeCactus())
	return &DesertBiome{
		structures: structures,
		perlin:       perlin.NewPerlin(2, 2, 1, seed),
		index: index,
	}
}

func makeCactus() *structure {
	s := makeStructure(1, 4, 1)
	s.blocks[0][0][0] = block.Cactus
	s.blocks[0][1][0] = block.Cactus
	s.blocks[0][2][0] = block.Cactus
	s.blocks[0][3][0] = block.Cactus
	s.p = 0.005
	return s
}

func (d *DesertBiome) blockType(x, y, z int, noises *noisesWithNeighbors) block.Block {
	height := noises.getNoise(x, z).elevation
	if y <= height{
		return block.Sand
	}
	return block.Air
}

func (d *DesertBiome) getScale() int {
	return desertElevationScale
}
func (d *DesertBiome) maxElevation() int {
	return desertMaxElevation
}
func (d *DesertBiome) minElevation() int {
	return desertMinElevation
}