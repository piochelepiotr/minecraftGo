package worldcontent

import (
	"github.com/aquilax/go-perlin"
	"github.com/piochelepiotr/minecraftGo/world/block"
)

const (
	forestMinElevation = 1
	// max is not reached, max - 1 is reached
	forestMaxElevation     = WorldHeight
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
	structures = append(structures, makeTree())
	return &ForestBiome{
		structures: structures,
		perlin:       perlin.NewPerlin(2, 2, 2, 233),
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
	return noise2d(f.perlin, x, z, forestScale, forestMinElevation, forestMaxElevation)
}
