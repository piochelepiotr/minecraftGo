package worldcontent

import (
	"github.com/aquilax/go-perlin"
	"github.com/piochelepiotr/minecraftGo/random"
	"github.com/piochelepiotr/minecraftGo/world/block"
)

const (
	plainTreeProbability float64 = 0.05
	plainTallGrassProbability float64 = 0.02
	plainMinElevation = 60
	// max is not reached, max - 1 is reached
	plainMaxElevation     = 65
	plainScale            = 100
	plainDirtLayerThikness  int = 5
)

func makePlainTree() *structure {
	s := makeStructure(5, 7, 5)
	for x := 0; x < 5; x++ {
		for y := 3; y < 5; y++ {
			for z := 0; z < 5; z++ {
				s.blocks[x][y][z] = block.BirchLeaves
			}
		}
	}
	for x := 1; x < 4; x++ {
		for z := 1; z < 4; z++ {
			s.blocks[x][5][z] = block.BirchLeaves
		}
	}
	for y := 0; y < 6; y++ {
		s.blocks[2][y][2] = block.Birch
	}
	for i := 1; i < 4; i++ {
		s.blocks[i][6][2] = block.BirchLeaves
		s.blocks[2][6][i] = block.BirchLeaves
	}
	s.p = plainTreeProbability
	s.originX = 2
	s.originZ = 2
	return s
}

func makePlainTallGrass() *structure {
	s := makeStructure(1, 1, 1)
	s.blocks[0][0][0] = block.TallGrass
	s.p = plainTallGrassProbability
	return s
}

type PlainBiome struct {
	index int
	structures []*structure
	perlin     *perlin.Perlin
	noise *random.Noise
}

func (f *PlainBiome) getStructures() []*structure {
	return f.structures
}

func makePlainBiome(seed int64, index int) *PlainBiome {
	plainSeed := seed * 3
	structures := make([]*structure, 0)
	structures = append(structures, makePlainTree())
	structures = append(structures, makePlainTallGrass())
	return &PlainBiome{
		structures: structures,
		perlin:     perlin.NewPerlin(2, 2, 3, plainSeed),
		noise: random.NewNoise(plainSeed),
		index: index,
	}
}

func (f *PlainBiome) blockType(x, y, z int, distanceFromBorder float64, noises *noisesWithNeighbors) block.Block {
	if y >= WorldHeight {
		return block.Air
	}
	if y == 0 {
		return block.BedRock
	}
	height := f.worldHeight(x, z, distanceFromBorder, noises)
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

func (f *PlainBiome) worldHeight(x, z int, distanceFromBorder float64, noises *noisesWithNeighbors) int {
	return elevation(noises.getNoise(x, z).elevationNoises[f.index], plainMinElevation, plainMaxElevation, distanceFromBorder)
}
