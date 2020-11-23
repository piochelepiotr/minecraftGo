package worldcontent

import (
	"github.com/piochelepiotr/minecraftGo/world/block"
)

const (
	plainMinElevation = 60
	// max is not reached, max - 1 is reached
	plainMaxElevation     = 65
	plainScale            = 100
	plainDirtLayerThikness  int = 5
)

func plainBlockType(y, elevation int, oreNoise, caveNoise float64) block.Block {
	b := block.Stone
	if y == elevation {
		b = block.Grass
	} else if y > elevation-plainDirtLayerThikness {
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