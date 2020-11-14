package worldcontent

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/piochelepiotr/minecraftGo/geometry"
	"github.com/piochelepiotr/minecraftGo/world/block"
	"log"
	"math"
	"sync"
)

const (
	UIDeleteChunkDistance  float32 = 120
	UILoadChunkDistance    float32 = 100
	deleteChunkDistance float32 = UIDeleteChunkDistance + 20
	maxWallJump float32 = 0.4
	backwardJump float32 = 0.4
)

// InMemoryWorld is a cache to the world on disk / generated by the Generator
type InMemoryWorld struct {
	// this is not thread safe for now
	chunks       map[geometry.Point]RawChunk
	chunksLock sync.RWMutex
	generator *Generator
	config Config
	chunksToWrite chan RawChunk
	cacheMisses int
}

func (w *InMemoryWorld) Config() Config {
	return w.config
}

func (w *InMemoryWorld) OutChunksToWrite() <-chan RawChunk {
	return w.chunksToWrite
}

func NewInMemoryWorld(config Config) *InMemoryWorld{
	return &InMemoryWorld{
		chunks: make(map[geometry.Point]RawChunk),
		config: config,
		generator: newGenerator(config),
		chunksToWrite: make(chan RawChunk, 200),
	}
}

func (w *InMemoryWorld) getChunk(p geometry.Point) RawChunk{
	w.chunksLock.Lock()
	defer w.chunksLock.Unlock()
	if chunk, ok := w.chunks[p]; ok {
		return chunk
	}
	w.cacheMisses++
	chunk := getChunk(w.config,  p, w.generator)
	w.chunks[p] = chunk
	return chunk
}

func ChunkStart(x int) int {
	return int(math.Floor(float64(x)/float64(ChunkSize))) * ChunkSize
}

func (w *InMemoryWorld) GetBlock(x, y, z int) block.Block {
	chunkX:= ChunkStart(x)
	chunkY := ChunkStart(y)
	chunkZ := ChunkStart(z)
	if y < 0 {
		return block.Stone
	}
	if y >= WorldHeight {
		return block.Air
	}
	c := w.getChunk(geometry.Point{chunkX, chunkY, chunkZ})
	return c.GetBlock(x - chunkX, y - chunkY, z - chunkZ)
}

// GetHeight returns height of the world in blocks at a x,z position
func (w *InMemoryWorld) GetHeight(x, z int) int {
	for y := WorldHeight - 1; y >= 0; y-- {
		if w.GetBlock(x, y, z) != block.Air {
			return y + 1
		}
	}
	return 0
}

//SetBlock sets a block and update the chunk
func (w *InMemoryWorld) SetBlock(x, y, z int, b block.Block) {
	w.chunksLock.Lock()
	defer w.chunksLock.Unlock()
	chunkX := ChunkStart(x)
	chunkY := ChunkStart(y)
	chunkZ := ChunkStart(z)
	p := geometry.Point{
		X: chunkX,
		Y: chunkY,
		Z: chunkZ,
	}
	if chunk, ok := w.chunks[p]; ok {
		chunk.SetBlock(x-chunkX, y-chunkY, z-chunkZ, b)
	} else {
		log.Print("ERROR when setting block in chunk ", p, " chunk isn't loaded")
	}
}

// even if the point is a bit inside a wall, this is going to return
func (w *InMemoryWorld) PlaceInFrontWithJumps(edges []mgl32.Vec3, dir mgl32.Vec3) float32 {
	min := float32(1000)
	for _, p := range edges {
		place := w.PlaceInFrontWithJumpsOnePoint(p, dir)
		if place < min {
			min = place
		}
	}
	return min
}

// even if the point is a bit inside a wall, this is going to return
func (w *InMemoryWorld) PlaceInFrontWithJumpsOnePoint(p mgl32.Vec3, dir mgl32.Vec3) float32 {
	place, _, _ := w.PlaceInFront(p, dir)
	if place > 0 {
		return place
	}
	uDir := dir.Mul(1/dir.Len())
	placeWithJump, _, _ := w.PlaceInFront(p.Add(uDir.Mul(maxWallJump)), dir)
	if placeWithJump > 0 {
		return placeWithJump + maxWallJump
	}
	if dir.Y() == 0  && (dir.X() == 0 || dir.Z() == 0){
		orthDir := mgl32.Vec3{uDir.Z(), 0, uDir.X()}
		placeWithBackJump, _, _ := w.PlaceInFront(p.Add(orthDir.Mul(backwardJump)), dir)
		if placeWithBackJump > 0 {
			return placeWithBackJump
		}
		placeWithForwardJump, _, _ := w.PlaceInFront(p.Add(orthDir.Mul(-backwardJump)), dir)
		if placeWithForwardJump > 0 {
			return placeWithForwardJump
		}
	}
	return 0
}

//PlaceInFront returns place in front of the player
func (w *InMemoryWorld) PlaceInFront(p mgl32.Vec3, dir mgl32.Vec3) (float32, geometry.Point, geometry.Point) {
	xray := geometry.NewXray(p, dir)
	for i := 0; i < 10; i++ {
		//fmt.Println("POS: ", x, ";", y, ";", z)
		if w.GetBlock(xray.P.X, xray.P.Y, xray.P.Z) != block.Air {
			return xray.Distance, xray.P, xray.Previous
		}
		xray.GoToNextBlock()
	}
	return xray.Distance, geometry.Point{}, geometry.Point{}
}

func (w *InMemoryWorld) ExpireChunks(playerPos mgl32.Vec3) {
	w.chunksLock.Lock()
	defer w.chunksLock.Unlock()
	for p, chunk := range w.chunks {
		if chunk.start.DistanceTo(playerPos) > deleteChunkDistance {
			if chunk.dirty {
				w.chunksToWrite <- chunk
			}
			delete(w.chunks, p)
		}
	}
}

// Close saves the world when closing the game
func (w *InMemoryWorld) Close() {
	for _, chunk := range w.chunks {
		if chunk.dirty {
			w.chunksToWrite <- chunk
		}
	}
	close(w.chunksToWrite)
	log.Println("Closed in memory world world")
}