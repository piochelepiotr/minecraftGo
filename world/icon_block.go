package world

import (
	"github.com/piochelepiotr/minecraftGo/game_engine/loader"
	"github.com/piochelepiotr/minecraftGo/game_engine/models"
	"github.com/piochelepiotr/minecraftGo/world/block"
)

// this should be moved to another package
// Change origin is a hack for the inventory
func (c *Chunk) ChangeOrigin() {
	// now the origin of the cube is in the middle
	for i := range c.model.vertices {
		c.model.vertices[i][0] -= 0.5
		c.model.vertices[i][1] -= 0.5
		c.model.vertices[i][2] -= 0.5
	}
}

func GetIconBlock(b block.Block, loader *loader.Loader) models.RawModel {
	c := newConstructionChunk()
	c.buildBlock(-0.5, -0.5, -0.5, b, true, true, true, true, true, true)
	model := loader.LoadToVAO(flatten3D(c.vertices), flatten2D(c.textures), c.indices, flatten3D(c.normals), flatten3D(c.colors))
	return model
}