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
	UIDeleteChunkDistance  float32 = 220
	UILoadChunkDistance    float32 = 200
	deleteChunkDistance float32 = UIDeleteChunkDistance + 20
	maxWallJump float32 = 0.4
	backwardJump float32 = 0.4
)

type blockUpdate struct {
	p geometry.Point
	b block.Block
}

// InMemoryWorld is a cache to the world on disk / generated by the Generator
type InMemoryWorld struct {
	// this is not thread safe for now
	chunks       map[geometry.Point]*RawChunk
	chunksLock sync.RWMutex
	generator *Generator
	config Config
	chunksToWrite chan *RawChunk
	cacheMisses int

	// used for struct editor
	structEditing bool
	blockUpdates []blockUpdate
}

func (w *InMemoryWorld) Config() Config {
	return w.config
}

func (w *InMemoryWorld) OutChunksToWrite() <-chan *RawChunk {
	return w.chunksToWrite
}

func NewInMemoryWorld(config Config, structEditing bool) *InMemoryWorld{
	return &InMemoryWorld{
		chunks: make(map[geometry.Point]*RawChunk),
		config: config,
		generator: newGenerator(config),
		chunksToWrite: make(chan *RawChunk, 200),
		structEditing: true,
	}
}

func (w *InMemoryWorld) GetChunk(p geometry.Point) *RawChunk{
	w.chunksLock.Lock()
	defer w.chunksLock.Unlock()
	if chunk, ok := w.chunks[p]; ok {
		return chunk
	}
	w.cacheMisses++
	// load blocks column by column
	chunkColumn := getChunkColumn(w.config, geometry.Point2D{p.X, p.Z}, w.generator)
	for y, chunk := range chunkColumn {
		p = geometry.Point{p.X, y*ChunkSize, p.Z}
		if _, ok := w.chunks[p]; !ok {
			w.chunks[p] = chunk
		}
	}
	return w.chunks[p]
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
	c := w.GetChunk(geometry.Point{chunkX, chunkY, chunkZ})
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
func (w *InMemoryWorld) SetBlock(x, y, z int, b block.Block) (updated bool){
	w.chunksLock.Lock()
	defer w.chunksLock.Unlock()
	if w.structEditing {
		w.blockUpdates = append(w.blockUpdates, blockUpdate{p: geometry.Point{x, y, z}, b: b})
	}
	chunkX := ChunkStart(x)
	chunkY := ChunkStart(y)
	chunkZ := ChunkStart(z)
	p := geometry.Point{
		X: chunkX,
		Y: chunkY,
		Z: chunkZ,
	}
	if chunk, ok := w.chunks[p]; ok {
		return chunk.SetBlock(x-chunkX, y-chunkY, z-chunkZ, b)
	}
	log.Print("ERROR when setting block in chunk ", p, " chunk isn't loaded")
	return false
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
	place, _, _ := w.GetPointedBlock(p, dir, true)
	if place > 0 {
		return place
	}
	uDir := dir.Mul(1/dir.Len())
	placeWithJump, _, _ := w.GetPointedBlock(p.Add(uDir.Mul(maxWallJump)), dir, true)
	if placeWithJump > 0 {
		return placeWithJump + maxWallJump
	}
	if dir.Y() == 0  && (dir.X() == 0 || dir.Z() == 0){
		orthDir := mgl32.Vec3{uDir.Z(), 0, uDir.X()}
		placeWithBackJump, _, _ := w.GetPointedBlock(p.Add(orthDir.Mul(backwardJump)), dir, true)
		if placeWithBackJump > 0 {
			return placeWithBackJump
		}
		placeWithForwardJump, _, _ := w.GetPointedBlock(p.Add(orthDir.Mul(-backwardJump)), dir, true)
		if placeWithForwardJump > 0 {
			return placeWithForwardJump
		}
	}
	return 0
}

//GetPointedBlock returns place in front of the player
func (w *InMemoryWorld) GetPointedBlock(p mgl32.Vec3, dir mgl32.Vec3, solid bool) (float32, geometry.Point, geometry.Point) {
	xray := geometry.NewXray(p, dir)
	for i := 0; i < 10; i++ {
		//fmt.Println("POS: ", x, ";", y, ";", z)
		b := w.GetBlock(xray.P.X, xray.P.Y, xray.P.Z)
		if (solid && b.IsSolid()) || (!solid && b != block.Air){
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

func (w *InMemoryWorld) saveStruct() {
	if len(w.blockUpdates) == 0 {
		return
	}
	p0 := w.blockUpdates[0].p
	minX := p0.X
	minY := p0.Y
	minZ := p0.Z
	maxX := p0.X
	maxY := p0.Y
	maxZ := p0.Z
	for _, u := range w.blockUpdates {
		p := u.p
		if p.X < minX {
			minX = p.X
		}
		if p.Y < minY {
			minY = p.Y
		}
		if p.Z < minZ {
			minZ = p.Z
		}
		if p.X > maxX {
			maxX = p.X
		}
		if p.Y > maxY {
			maxY = p.Y
		}
		if p.Z > maxZ {
			maxZ = p.Z
		}
	}
	s := savedStructure{
		X : maxX - minX + 1,
		Y : maxY - minY + 1,
		Z : maxZ - minZ + 1,
	}
	s.Blocks = make([]byte, s.X * s.Y * s.Z)
	for i := 0; i < s.X * s.Y * s.Z; i++ {
		s.Blocks[i] = byte(block.Air)
	}
	for _, u := range w.blockUpdates {
		s.Blocks[structIndex(s.Y, s.Z, u.p.X-minX, u.p.Y-minY, u.p.Z-minZ)] = byte(u.b)
	}
	err := s.save("struct")
	if err != nil {
		log.Printf("Problem saving structure")
		return
	}
	log.Printf("saved structure")
}

// Close saves the world when closing the game
func (w *InMemoryWorld) Close() {
	if w.structEditing {
		w.saveStruct()
	}

	for _, chunk := range w.chunks {
		if chunk.dirty {
			w.chunksToWrite <- chunk
		}
	}
	close(w.chunksToWrite)
	log.Println("Closed in memory world world")
}
