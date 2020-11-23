package worldcontent

import (
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

func forestBlockType(y, elevation int, oreNoise, caveNoise float64) block.Block {
	b := block.Stone
	if y == elevation {
		b = block.Grass
	} else if y > elevation-dirtLayerThikness {
		b = block.Dirt
	}
	// y and z inverted on purpose, negative 3rd metric doesn't work
	if caveNoise < 0.1 {
		return block.Air
	}
	if b != block.Stone {
		return b
	}
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