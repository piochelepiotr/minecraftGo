package world

import (
	"github.com/piochelepiotr/minecraftGo/geometry"
	"io/ioutil"
	"os"
)

type ChunkLoader struct {
	LoadedChunk chan *Chunk
	generator *Generator
	worldConfig Config
}

func NewChunkLoader(worldConfig Config, generator *Generator) *ChunkLoader {
	return &ChunkLoader{
		worldConfig: worldConfig,
		generator: generator,
		LoadedChunk: make(chan *Chunk, 100),
	}
}

func (c *ChunkLoader) Run(chunkLoadDecisions <-chan geometry.Point) {
	go func() {
		for p := range chunkLoadDecisions {
			chunk := GetGraphicChunk(c.worldConfig, p, c.generator)
			c.LoadedChunk <- chunk
		}
		close(c.LoadedChunk)
	}()
}

// GetChunk tries to load the chunk from a saved file. If there is nothing, generates one using the generator
func GetChunk(worldConfig Config, start geometry.Point, generator *Generator) RawChunk {
	if chunk, err := LoadChunkFromSaves(worldConfig, start); err == nil {
		return chunk
	}
	return generator.GenerateChunk(start)
}

func GetGraphicChunk(worldConfig Config, start geometry.Point, generator *Generator) *Chunk {
	return NewChunk(GetChunk(worldConfig, start, generator))
}

func LoadChunkFromSaves(worldConfig Config, start geometry.Point) (chunk RawChunk, err error) {
	path := savesPath + "/" + worldConfig.Name + "/" + "chunk_" + start.GetKey()
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
