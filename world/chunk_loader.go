package world

import (
	"github.com/piochelepiotr/minecraftGo/geometry"
	"io/ioutil"
	"os"
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
			chunk := GetGraphicChunk(p, c.generator)
			c.LoadedChunk <- chunk
		}
		close(c.LoadedChunk)
	}()
}

// GetChunk tries to load the chunk from a saved file. If there is nothing, generates one using the generator
func GetChunk(start geometry.Point, generator *Generator) RawChunk {
	if chunk, err := LoadChunkFromSaves(start); err == nil {
		return chunk
	}
	return generator.GenerateChunk(start)
}

func GetGraphicChunk(start geometry.Point, generator *Generator) *Chunk {
	chunk := Chunk{RawChunk: GetChunk(start, generator)}
	chunk.buildFaces()
	return &chunk
}

func LoadChunkFromSaves(start geometry.Point) (chunk RawChunk, err error) {
	path := savesPath + "/" + "chunk_" + start.GetKey()
	if _, err := os.Stat(path); err != nil {
		return chunk, err
	}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return chunk, err
	}
	chunk = decode(data, start)
	return chunk, nil
}
