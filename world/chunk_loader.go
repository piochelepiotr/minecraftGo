package world

import (
	"github.com/piochelepiotr/minecraftGo/geometry"
)

type ChunkLoader struct {
	LoadedChunk chan *Chunk
	generator *Generator
}

func NewChunkLoader(generator *Generator) *ChunkLoader {
	return &ChunkLoader{
		generator: generator,
		LoadedChunk: make(chan *Chunk, 100),
	}
}

func (c *ChunkLoader) Run(chunkLoadDecisions <-chan geometry.Point) {
	go func() {
		for p := range chunkLoadDecisions {
			chunk := c.generator.GenerateChunk(p)
			c.LoadedChunk <- chunk
		}
		close(c.LoadedChunk)
	}()
}
