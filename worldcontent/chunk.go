package worldcontent

import (
	"github.com/piochelepiotr/minecraftGo/geometry"
	"github.com/piochelepiotr/minecraftGo/world/block"
)

const (
	chunkFormatV1 = 1
)

type RawChunk struct {
	blocks []block.Block
	start            geometry.Point
	// dirty is true when the content of the chunk hasn't been save to disk yet
	dirty bool
}

func index(x, y, z int) int {
	return x*ChunkSize*ChunkSize+y*ChunkSize+z
}

// SetBlock sets a block in a chunk
func (c *RawChunk) SetBlock(x, y, z int, b block.Block) {
	c.blocks[index(x, y, z)] = b
	c.dirty = true
}

// GetBlock gets the block of a chunk
func (c *RawChunk) GetBlock(x, y, z int) block.Block {
	return c.blocks[index(x, y, z)]
}

// GetHeight gets the height of the chunk in blocks (not including air)
func (c *RawChunk) GetHeight(x, z int) int {
	for y := ChunkSize - 1; y >= 0; y-- {
		if c.GetBlock(x, y, z) != block.Air {
			return y + 1
		}
	}
	return 0
}
