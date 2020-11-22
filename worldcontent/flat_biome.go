package worldcontent

import "github.com/piochelepiotr/minecraftGo/world/block"

func flatBlockType(y, elevation int, oreNoise, caveNoise float64) block.Block {
	if y == elevation {
		return block.Grass
	} else if y > elevation-plainDirtLayerThikness {
		return block.Dirt
	}
	return block.Stone
}
