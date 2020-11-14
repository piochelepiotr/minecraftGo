package world

import (
	"github.com/piochelepiotr/minecraftGo/geometry"
	"github.com/piochelepiotr/minecraftGo/worldcontent"
)

// ChunkLoader loads the chunks in parallel to the main game happening
type ChunkLoader struct {
	chunks chan *Chunk
	world *worldcontent.InMemoryWorld
}

func NewChunkLoader(world *worldcontent.InMemoryWorld, chunksToLoad <-chan geometry.Point) *ChunkLoader {
	l := &ChunkLoader{
		world: world,
		chunks: make(chan *Chunk, 100),
	}
	l.start(chunksToLoad)
	return l
}

func (c *ChunkLoader) Chunks() <-chan *Chunk {
	return c.chunks
}

func (c *ChunkLoader) start(chunksToLoad <-chan geometry.Point) {
	// for now, only concurrency of 1
	go func() {
		for p := range chunksToLoad {
			chunk := c.GetChunk(p)
			c.chunks <- chunk
		}
		close(c.chunks)
	}()
}

func (c *ChunkLoader) GetChunk(p geometry.Point) *Chunk {
	return NewChunk(c.world, p)
}
