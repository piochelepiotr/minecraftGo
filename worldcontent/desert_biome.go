package worldcontent

import (
	"github.com/piochelepiotr/minecraftGo/perlin"
	"github.com/piochelepiotr/minecraftGo/world/block"
)

const (
	desertMinElevation = 1
	// max is not reached, max - 1 is reached
	desertMaxElevation = 20
	desertElevationScale = 100
)

type DesertBiome struct {
	structures []*structure
	perlin *perlin.Perlin
}

func (d *DesertBiome) getStructures() []*structure {
	return d.structures
}

func makeDesertBiome() *DesertBiome {
	structures := make([]*structure, 0)
	structures = append(structures, makeCactus())
	return &DesertBiome{
		structures: structures,
		perlin:       perlin.NewPerlin(345),
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

func (d *DesertBiome) blockType(x, y, z int) block.Block {
	height := d.worldHeight(x, z)
	if y <= height{
		return block.Sand
	}
	return block.Air
}

func (d *DesertBiome) worldHeight(x, z int) int {
	c := d.perlin.Noise2D(float64(x)/desertElevationScale, float64(z)/desertElevationScale)
	return desertMinElevation + int(float64(desertMaxElevation - desertMinElevation)*c)
}
