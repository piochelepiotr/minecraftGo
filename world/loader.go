package world

import (
	"github.com/piochelepiotr/minecraftGo/geometry"
	"github.com/piochelepiotr/minecraftGo/worldcontent"
	"log"
	"time"
)

// ChunkLoader loads the chunks in parallel to the main game happening
type ChunkLoader struct {
	chunks chan *Chunk
	world *worldcontent.InMemoryWorld
	loaded int
}

func NewChunkLoader(world *worldcontent.InMemoryWorld, chunksToLoad <-chan geometry.Point) *ChunkLoader {
	l := &ChunkLoader{
		world: world,
		chunks: make(chan *Chunk, 100),
	}
	go l.start(chunksToLoad)
	return l
}

func (c *ChunkLoader) Chunks() <-chan *Chunk {
	return c.chunks
}

func (c *ChunkLoader) start(chunksToLoad <-chan geometry.Point) {
	// for now, only concurrency of 1
	report := time.NewTicker(time.Second*2)
	defer close(c.chunks)
	for {
		select {
		case p, ok := <-chunksToLoad:
			if !ok {
				return
			}
			chunk := c.GetChunk(p)
			c.chunks <- chunk
			c.loaded++
		case <-report.C:
			log.Println("chunks loaded per s", c.loaded/2)
			c.loaded = 0
		}
	}
}

func (c *ChunkLoader) GetChunk(p geometry.Point) *Chunk {
	return NewChunk(c.world, p)
}
