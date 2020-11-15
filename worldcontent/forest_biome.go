package worldcontent

import (
	"github.com/piochelepiotr/minecraftGo/perlin"
	"github.com/piochelepiotr/minecraftGo/world/block"
)

const (
	forestMinElevation = 1
	// max is not reached, max - 1 is reached
	forestMaxElevation     = 16 * 3
	forestScale            = 100
	dirtLayerThikness  int = 5
)


type ForestBiome struct {
	structures []*structure
	perlin *perlin.Perlin
}

func (f *ForestBiome) getStructures() []*structure {
	return f.structures
}

func makeForestBiome() *ForestBiome {
	structures := make([]*structure, 0)
	// structures = append(structures, makeTree())
	return &ForestBiome{
		structures: structures,
		perlin:       perlin.NewPerlin(233),
	}
}

func (f *ForestBiome) blockType(x, y, z int) block.Block {
	height := f.worldHeight(x, z)
	if y == height {
		return block.Grass
	}
	if y < height && y > height - dirtLayerThikness {
		return block.Dirt
	}
	if y <= height - dirtLayerThikness {
		return block.Stone
	}
	return block.Air
}

func (f *ForestBiome) worldHeight(x, z int) int {
	c := f.perlin.Noise2D(float64(x)/float64(forestScale), float64(z)/float64(forestScale))
	return forestMinElevation + int(float64(forestMaxElevation - forestMinElevation)*c)
}
