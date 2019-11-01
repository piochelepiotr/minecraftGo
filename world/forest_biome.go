package world

import (
	"github.com/aquilax/go-perlin"
)

const (
	forestMinElevation = 1
	// max is not reached, max - 1 is reached
	forestMaxElevation = WorldHeight
	forestElevationScale = 100
	dirtLayerThikness int = 5
)


type ForestBiome struct {
	structures []*Structure
	perlin *perlin.Perlin
}

func (f *ForestBiome) getStructures() []*Structure {
	return f.structures
}

func makeForestBiome() *ForestBiome {
	structures := make([]*Structure, 0)
	structures = append(structures, makeTree())
	return &ForestBiome{
		structures: structures,
		perlin:       perlin.NewPerlin(alpha, beta, perlinN, 233),
	}
}

func (f *ForestBiome) blockType(x, y, z int) Block {
	height := f.worldHeight(x, z)
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

func (f *ForestBiome) worldHeight(x, z int) int {
	c := perlinCoef(f.perlin, x, z, forestElevationScale)
	return forestMinElevation + int(float64(forestMaxElevation - forestMinElevation)*c)
}
