package worldcontent

import (
	"github.com/aquilax/go-perlin"
	"github.com/piochelepiotr/minecraftGo/random"
	"github.com/piochelepiotr/minecraftGo/world/block"
)

const (
	treeProbability float64 = 0.01
	tallGrassProbability float64 = 0.1
	birchSaplingProbability float64 = 0.05
	roseProbability float64 = 0.05
	forestMinElevation = 60
	// max is not reached, max - 1 is reached
	forestMaxElevation     = 100
	forestScale            = 50
	dirtLayerThikness  int = 5
	goldProba float64 = 0.01
	ironProba float64 = 0.05
	coalProba float64 = 0.1
)

func makeTree() *structure {
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
	s.p = treeProbability
	s.originX = 2
	s.originZ = 2
	return s
}

func makeTallGrass() *structure {
	s := makeStructure(1, 1, 1)
	s.blocks[0][0][0] = block.TallGrass
	s.p = tallGrassProbability
	return s
}

func makeBirchSampling() *structure {
	s := makeStructure(1, 1, 1)
	s.blocks[0][0][0] = block.BirchSapling
	s.p = birchSaplingProbability
	return s
}

func makeRose() *structure {
	s := makeStructure(1, 1, 1)
	s.blocks[0][0][0] = block.Rose
	s.p = roseProbability
	return s
}

type ForestBiome struct {
	structures []*structure
	perlin     *perlin.Perlin
	noise *random.Noise
	index int
}

func (f *ForestBiome) getStructures() []*structure {
	return f.structures
}

func makeForestBiome(seed int64, index int) *ForestBiome {
	forestSeed := seed * 2
	structures := make([]*structure, 0)
	structures = append(structures, makeTree())
	structures = append(structures, makeTallGrass())
	structures = append(structures, makeBirchSampling())
	structures = append(structures, makeRose())
	return &ForestBiome{
		structures: structures,
		perlin:     perlin.NewPerlin(1.3, 2, 3, forestSeed),
		noise: random.NewNoise(forestSeed),
		index: index,
	}
}

func (f *ForestBiome) blockType(x, y, z int, noises *noisesWithNeighbors) block.Block {
	if y >= WorldHeight {
		return block.Air
	}
	if y == 0 {
		return block.BedRock
	}
	height := noises.getNoise(x, z).elevation
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

func (f *ForestBiome) getScale() int {
	return forestScale
}
func (f *ForestBiome) maxElevation() int {
	return forestMaxElevation
}
func (f *ForestBiome) minElevation() int {
	return forestMinElevation
}
