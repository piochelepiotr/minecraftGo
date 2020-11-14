package worldcontent

import (
	"github.com/piochelepiotr/minecraftGo/geometry"
	"github.com/piochelepiotr/minecraftGo/world/block"
)

const (
	chunkFormatV1 = 1
)


type RawChunk struct {
	size int
	blocks []block.Block
	start            geometry.Point
	// dirty is true when the content of the chunk hasn't been save to disk yet
	dirty bool
}

func NewOneBlockChunk(b block.Block) RawChunk{
	return RawChunk{
		size: 1,
		blocks: []block.Block{b},
	}
}

func (c *RawChunk) Size() int {
	return c.size
}

func (c *RawChunk) index(x, y, z int) int {
	return x*c.size*c.size+y*c.size+z
}

// SetBlock sets a block in a chunk
func (c *RawChunk) SetBlock(x, y, z int, b block.Block) {
	c.blocks[c.index(x, y, z)] = b
	c.dirty = true
}

// GetBlock gets the block of a chunk
func (c *RawChunk) GetBlock(x, y, z int) block.Block {
	return c.blocks[c.index(x, y, z)]
}

// GetHeight gets the height of the chunk in blocks (not including air)
func (c *RawChunk) GetHeight(x, z int) int {
	for y := c.size - 1; y >= 0; y-- {
		if c.GetBlock(x, y, z) != block.Air {
			return y + 1
		}
	}
	return 0
}
