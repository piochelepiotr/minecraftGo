package worldcontent

import (
	"github.com/piochelepiotr/minecraftGo/world/block"
)

const (
	desertMinElevation = 60
	// max is not reached, max - 1 is reached
	desertMaxElevation = 65
	desertScale        = 100
)

func desertBlockType(y, elevation int, oreNoise, caveNoise float64) block.Block {
	return block.Sand
}
