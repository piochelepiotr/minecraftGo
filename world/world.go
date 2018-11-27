package world

import (
	"math"

	"github.com/aquilax/go-perlin"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/piochelepiotr/minecraftGo/entities"
	"github.com/piochelepiotr/minecraftGo/models"
	"github.com/piochelepiotr/minecraftGo/textures"
)

// World contains all the blocks of the world in chunks that load around the player
type World struct {
	chunks       map[Point]Chunk
	modelTexture textures.ModelTexture
	perlin       *perlin.Perlin
}

func getChunk(x int) int {
	return int(math.Floor(float64(x) / float64(ChunkSize)))
}

// CreateWorld initiate the world
func CreateWorld(modelTexture textures.ModelTexture) World {
	chunks := make(map[Point]Chunk)
	return World{
		chunks:       chunks,
		modelTexture: modelTexture,
		perlin:       perlin.NewPerlin(alpha, beta, perlinN, 233),
	}
}

//GetChunks returns all chunks that are going to be rendered
func (w *World) GetChunks() []entities.Entity {
	chunks := make([]entities.Entity, 0)
	for p, chunk := range w.chunks {
		model := chunk.Model
		if model.VertexCount == 0 {
			continue
		}
		t := models.TexturedModel{
			RawModel:     chunk.Model,
			ModelTexture: w.modelTexture,
		}
		e := entities.Entity{
			TexturedModel: t,
			Position: mgl32.Vec3{
				float32(p.X),
				float32(p.Y),
				float32(p.Z),
			},
		}
		chunks = append(chunks, e)
	}
	return chunks
}

// LoadChunk loads a chunk into the world so that it is rendered
func (w *World) LoadChunk(x, y, z int) {
	p := Point{
		X: x,
		Y: y,
		Z: z,
	}
	w.chunks[p] = CreateChunk(x, y, z, w.modelTexture, w.perlin)
}

// GetHeight returns height of the world in blocks at a x,z position
func (w *World) GetHeight(x, z int) int {
	chunkX := getChunk(x)
	chunkZ := getChunk(z)
	for chunkY := WorldHeight - ChunkSize; chunkY >= 0; chunkY -= ChunkSize {
		if chunk, ok := w.chunks[Point{X: chunkX, Y: chunkY, Z: chunkZ}]; ok {
			return chunkY + chunk.GetHeight(x, z)
		}
	}
	return 0
}
