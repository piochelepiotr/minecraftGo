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

func makeCactus() *structure {
	s := makeStructure(1, 4, 1)
	s.blocks[0][0][0] = block.Cactus
	s.blocks[0][1][0] = block.Cactus
	s.blocks[0][2][0] = block.Cactus
	s.blocks[0][3][0] = block.Cactus
	s.p = 0.005
	return s
}

func desertBlockType(y, elevation int, oreNoise, caveNoise float64) block.Block {
	return block.Sand
}
