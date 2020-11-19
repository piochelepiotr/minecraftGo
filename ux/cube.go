package ux

import (
	"github.com/piochelepiotr/minecraftGo/game_engine/loader"
	"github.com/piochelepiotr/minecraftGo/game_engine/models"
)

func NewCubeOutline(loader *loader.Loader) models.OutlineModel {
	positions := []float32{}
	addPosition := func(x, y, z float32) {
		positions = append(positions, x)
		positions = append(positions, y)
		positions = append(positions, z)
	}
	for y := float32(0); y < 2; y++ {
		for x := float32(0); x < 2; x++ {
			for z := float32(0); z < 2; z++ {
				addPosition(x, y, z)
			}
		}
	}
	// bottom
	// ---> x
	// | 0 ---- 2
	// | |      |
	// z 1 ---- 3
	// up
	// ---> x
	// | 4 ---- 6
	// | |      |
	// z 5 ---- 7
	addPosition(0, 0, 0)
	addPosition(0, 0, 0)
	addPosition(0, 0, 0)
	indices := []uint32{
		0, 1, 2, 3, 0, 2, 1, 3,// bottom
		4, 5, 6, 7, 4, 6, 5, 7,// up
		0, 4, 1, 5, 2, 6, 3, 7,// connection between the two
		}
	return models.OutlineModel{
		RawModel: loader.LoadLinesToVAO(positions, indices),
	}
}