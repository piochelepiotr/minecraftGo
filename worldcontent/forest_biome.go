package worldcontent

import (
	"github.com/aquilax/go-perlin"
	"github.com/piochelepiotr/minecraftGo/random"
	"github.com/piochelepiotr/minecraftGo/world/block"
)

const (
	forestMinElevation = 60
	// max is not reached, max - 1 is reached
	forestMaxElevation     = 100
	forestScale            = 50
	dirtLayerThikness  int = 5
	goldProba float64 = 0.01
	ironProba float64 = 0.05
	coalProba float64 = 0.1
)

type ForestBiome struct {
	structures []*structure
	perlin     *perlin.Perlin
	noise *random.Noise
}

func (f *ForestBiome) getStructures() []*structure {
	return f.structures
}

func makeForestBiome(seed int64) *ForestBiome {
	forestSeed := seed * 2
	structures := make([]*structure, 0)
	structures = append(structures, makeTree())
	return &ForestBiome{
		structures: structures,
		perlin:     perlin.NewPerlin(1.3, 2, 3, forestSeed),
		noise: random.NewNoise(forestSeed),
	}
}

func (f *ForestBiome) blockType(x, y, z int) block.Block {
	if y >= WorldHeight {
		return block.Air
	}
	if y == 0 {
		return block.BedRock
	}
	height := f.worldHeight(x, z)
	if y > height {
		return block.Air
	}
	b := block.Stone
	if y == height {
		b = block.Grass
	} else if y > height-dirtLayerThikness {
		b = block.Dirt
	}
	// y and z inverted on purpose, negative 3rd metric doesn't work
	noise := noise3d(f.perlin, x, z, y, 40.0)
	if noise < 0.1 {
		return block.Air
	}
	if b != block.Stone {
		return b
	}
	oreNoise := f.noise.Noise3D(x, y, z)
	cumulative := float64(0)
	if oreNoise < cumulative + goldProba {
		return block.Gold
	}
	cumulative += goldProba
	if oreNoise < cumulative + ironProba {
		return block.Iron
	}
	cumulative += ironProba
	if oreNoise < cumulative + coalProba {
		return block.Coal
	}
	cumulative += coalProba
	return b
}

func (f *ForestBiome) worldHeight(x, z int) int {
	return elevation(f.perlin, x, z, forestScale, forestMinElevation, forestMaxElevation)
}
