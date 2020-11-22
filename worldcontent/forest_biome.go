package worldcontent

import (
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