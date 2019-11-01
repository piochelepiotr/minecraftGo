package world

import "github.com/aquilax/go-perlin"

const (
	desertMinElevation = 1
	// max is not reached, max - 1 is reached
	desertMaxElevation = 20
	desertElevationScale = 100
)

type DesertBiome struct {
	structures []*Structure
	perlin *perlin.Perlin
}

func (d *DesertBiome) getStructures() []*Structure {
	return d.structures
}

func makeDesertBiome() *DesertBiome {
	structures := make([]*Structure, 0)
	structures = append(structures, makeCactus())
	return &DesertBiome{
		structures: structures,
		perlin:       perlin.NewPerlin(alpha, beta, perlinN, 345),
	}
}

func makeCactus() *Structure {
	s := makeStructure(1, 4, 1)
	s.blocks[0][0][0] = Cactus
	s.blocks[0][1][0] = Cactus
	s.blocks[0][2][0] = Cactus
	s.blocks[0][3][0] = Cactus
	s.p = treeProbability
	return s
}

func (d *DesertBiome) blockType(x, y, z int) Block {
	height := d.worldHeight(x, z)
	if y <= height{
		return Sand
	}
	return Air
}

func (d *DesertBiome) worldHeight(x, z int) int {
	c := perlinCoef(d.perlin, x, z, desertElevationScale)
	return desertMinElevation + int(float64(desertMaxElevation - desertMinElevation)*c)
}
