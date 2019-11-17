package world

import (
	"fmt"
	"github.com/piochelepiotr/minecraftGo/geometry"
	"github.com/piochelepiotr/minecraftGo/loader"
	"time"
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
			loader.Debug = true
			fmt.Println("hello")
			fmt.Println(p.X)
			fmt.Println(p.Y)
			fmt.Println(p.Z)
			time.Sleep(time.Second)
			chunk := CreateChunk(p, c.generator)
			fmt.Println("hello 2")
			c.LoadedChunk <- chunk
		}
		close(c.LoadedChunk)
	}()
}
