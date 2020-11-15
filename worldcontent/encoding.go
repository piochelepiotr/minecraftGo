package worldcontent

import (
	"github.com/piochelepiotr/minecraftGo/geometry"
	"github.com/piochelepiotr/minecraftGo/world/block"
)

func (c *RawChunk) encode() []byte {
	encoded := make([]byte, 1, len(c.blocks) + 1)
	encoded[0] = chunkFormatV1
	for _, b := range c.blocks {
		encoded = append(encoded, byte(b))
	}
	return encoded
}

func decode(data []byte, start geometry.Point) (chunk *RawChunk) {
	// for now, we only have v1
	chunk = &RawChunk{}
	data = data[1:]
	chunk.start = start
	chunk.blocks = make([]block.Block, 0, len(data))
	for _, b := range data {
		chunk.blocks = append(chunk.blocks, block.Block(b))
	}
	return chunk
}
