package world
import (
    "github.com/piochelepiotr/minecraftGo/entities"
    "github.com/piochelepiotr/minecraftGo/textures"
    "github.com/piochelepiotr/minecraftGo/models"
	"github.com/go-gl/mathgl/mgl32"
    "math"
)

const WordHeight int = ChunkSize * 10

type World struct {
    chunks map[Point] Chunk
    modelTexture textures.ModelTexture
}

func getChunk(x int) int {
    return int(math.Floor(float64(x)/float64(ChunkSize)))
}

func CreateWorld(modelTexture textures.ModelTexture) World {
    chunks := make(map[Point] Chunk)
    return World{
        chunks: chunks,
        modelTexture: modelTexture,
    }
}

func (w *World)  GetChunks() []entities.Entity {
    chunks := make([]entities.Entity, 0)
    for p, chunk := range w.chunks {
        t := models.TexturedModel{
            RawModel: chunk.Model,
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

func (w *World) LoadChunk(x, y, z int) {
    p := Point{
        X: x,
        Y: y,
        Z: z,
    }
    w.chunks[p] = CreateChunk(x, y, z, w.modelTexture)
}

func (w *World) GetHeight(x, z int) int {
    chunkX := getChunk(x)
    chunkZ := getChunk(z)
    for chunkY := WordHeight - ChunkSize; chunkY >= 0; chunkY -= ChunkSize {
        if chunk, ok := w.chunks[Point{X:chunkX, Y:chunkY, Z:chunkZ}]; ok {
            return chunkY + chunk.GetHeight(x, z)
        }
    }
    return 0
}
