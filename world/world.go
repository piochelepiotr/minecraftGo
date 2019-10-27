package world

import (
	"fmt"
	"math"
	"time"

	"github.com/aquilax/go-perlin"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/piochelepiotr/minecraftGo/entities"
	"github.com/piochelepiotr/minecraftGo/models"
	"github.com/piochelepiotr/minecraftGo/textures"
	"github.com/piochelepiotr/minecraftGo/toolbox"
)

const (
	deleteChunkDistance  float32 = 130
	loadChunkDistance    float32 = 40
	chunkMaxLoadDistance int     = 10
)

// World contains all the blocks of the world in chunks that load around the player
type World struct {
	chunks       map[Point]*Chunk
	modelTexture textures.ModelTexture
	perlin       *perlin.Perlin
}

func getChunk(x int) int {
	return int(math.Floor(float64(x)/float64(ChunkSize))) * ChunkSize
}

// CreateWorld initiate the world
func CreateWorld(modelTexture textures.ModelTexture) World {
	chunks := make(map[Point]*Chunk)
	return World{
		chunks:       chunks,
		modelTexture: modelTexture,
		perlin:       perlin.NewPerlin(alpha, beta, perlinN, 233),
	}
}

//GetChunks returns all chunks that are going to be rendered
func (w *World) GetChunks() []entities.Entity {
	chunks := make([]entities.Entity, 0)
	for _, chunk := range w.chunks {
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
				float32(chunk.Start.X),
				float32(chunk.Start.Y),
				float32(chunk.Start.Z),
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
	chunk := CreateChunk(x, y, z, w.modelTexture, w.perlin)
	w.chunks[p] = &chunk
}

func (w *World) loadChunkIfNotLoaded(x, y, z int) bool {
	if _, ok := w.chunks[Point{x, y, z}]; !ok {
		w.LoadChunk(x, y, z)
		return true
	}
	return false
}

// GetHeight returns height of the world in blocks at a x,z position
func (w *World) GetHeight(x, z int) int {
	for y := WorldHeight - 1; y >= 0; y-- {
		if w.GetBlock(x, y, z) != Air {
			return y + 1
		}
	}
	return 0
}

//GetBlock returns block x,y,z
func (w *World) GetBlock(x, y, z int) Block {
	chunkX := getChunk(x)
	chunkY := getChunk(y)
	chunkZ := getChunk(z)
	p := Point{
		X: chunkX,
		Y: chunkY,
		Z: chunkZ,
	}
	if chunk, ok := w.chunks[p]; ok {
		return chunk.GetBlock(x-chunkX, y-chunkY, z-chunkZ)
	}
	//fmt.Println("ERROR when getting block in chunk ", p)
	return Air
}

//SetBlock sets a block and update the chunk
func (w *World) SetBlock(x, y, z int, b Block) {
	chunkX := getChunk(x)
	chunkY := getChunk(y)
	chunkZ := getChunk(z)
	p := Point{
		X: chunkX,
		Y: chunkY,
		Z: chunkZ,
	}
	if chunk, ok := w.chunks[p]; ok {
		chunk.SetBlock(x-chunkX, y-chunkY, z-chunkZ, Air)
	} else {
		//fmt.Println("ERROR when setting block in chunk ", p)
	}
}

//PlaceInFront returns place in front of the player
func (w *World) PlaceInFront(px, py, pz float32, dir mgl32.Vec3) (float32, Point) {
	dist := float32(0)
	x := int(math.Floor(float64(px)))
	y := int(math.Floor(float64(py)))
	z := int(math.Floor(float64(pz)))
	for i := 0; i < 10; i++ {
		//fmt.Println("POS: ", x, ";", y, ";", z)
		if w.GetBlock(x, y, z) != Air {
			return dist, Point{x, y, z}
		}
		dist += toolbox.GetNextBlock(&px, &py, &pz, dir, &x, &y, &z)
	}
	return dist, Point{0, 0, 0}
}

func truncateMovement(move mgl32.Vec3, place float32) mgl32.Vec3 {
	if move.Len() > place {
		return move.Mul(place / move.Len())
	}
	// mul * 1 to make a copy
	return move.Mul(1)
}

// TouchesGround returns true if the player touches the ground
func (w *World) TouchesGround(player *entities.Player) bool {
	p := player.Entity.Position
	place, _ := w.PlaceInFront(p.X(), p.Y(), p.Z(), mgl32.Vec3{0, -1, 0})
	return place == 0
}

// MovePlayer moves the player inside the world
func (w *World) MovePlayer(player *entities.Player, forward, backward, jump, touchGround bool) {
	now := time.Now()
	if !player.LastMove.IsZero() {
		dt := float32(float64(now.Sub(player.LastMove)) / float64(time.Second))
		p := player.Entity.Position
		// vertical movement first
		if player.Speed.Y() != 0 {
			y := player.Speed.Y() * dt
			dir := float32(-1)
			if y > 0 {
				dir = 1
			}
			place, _ := w.PlaceInFront(p.X(), p.Y(), p.Z(), mgl32.Vec3{0, dir, 0})
			fmt.Printf("Speed y is %f\n", player.Speed.Y())
			if toolbox.Abs(y) > place {
				player.Speed = mgl32.Vec3{player.Speed.X(), 0, player.Speed.Z()}
				player.Entity.Position = mgl32.Vec3{p.X(), float32(math.Floor(float64(dir*place+player.Entity.Position.Y()))), p.Z()}
			} else {
				player.Entity.Position = player.Entity.Position.Add(mgl32.Vec3{0, y, 0})
			}
			fmt.Printf("Pos y is %f\n", player.Entity.Position.Y())
		}

		hSpeed := mgl32.Vec3{player.Speed.X(), 0, player.Speed.Z()}

		// horizontal movement second
		if hSpeed.Len() > 0 {
			move := hSpeed.Mul(dt)
			// fmt.Printf("Move is %f\n", move.Len())
			//forward := mgl32.Vec3{0, -1, 0}
			forward := hSpeed
			// first, go as far as we can in the forward direction
			place, _ := w.PlaceInFront(p.X(), p.Y(), p.Z(), forward)
			firstMove := truncateMovement(move, place)
			// fmt.Printf("first move is %f\n", firstMove.Len())
			player.Entity.Position = player.Entity.Position.Add(firstMove)
			restMove := move.Sub(firstMove)
			if restMove.Len() > 0 {
				if restMove.X() != 0 {
				}
			}
		}
		if !(forward || backward) || !touchGround {
			player.Speed = entities.Friction(player.Speed, dt)
		}
		player.Speed = entities.Gravity(player.Speed, dt, touchGround)
	}
	player.LastMove = now
}

// ClickOnBlock removes or adds a block
func (w *World) ClickOnBlock(camera *entities.Camera) {
	xray := toolbox.ComputeCameraRay(camera.Rotation)
	p := camera.Position
	place, block := w.PlaceInFront(p.X(), p.Y(), p.Z(), xray)
	fmt.Println(place)
	w.SetBlock(block.X, block.Y, block.Z, Air)
}

func (w *World) deleteChunks(playerPos mgl32.Vec3) {
	for p, chunk := range w.chunks {
		if chunk.Start.DistanceTo(playerPos) > deleteChunkDistance {
			delete(w.chunks, p)
		}
	}
}

//LoadAllChunks loads one chunk per second
func (w *World) LoadAllChunks(playerPos mgl32.Vec3) {
	for {
		time.Sleep(1e9)
		w.LoadChunks(playerPos)
	}
}

//LoadChunks load chunks around the player, for now, only load one chunk per frame
func (w *World) LoadChunks(playerPos mgl32.Vec3) {
	xPlayer := int(playerPos.X())
	zPlayer := int(playerPos.Z())
	chunkX := getChunk(xPlayer)
	chunkZ := getChunk(zPlayer)
	for i := 0; i < chunkMaxLoadDistance; i++ {
		p := Point{i * ChunkSize, 0, i * ChunkSize}
		if p.DistanceTo(playerPos) > loadChunkDistance {
			return
		}
		z := -i
		for x := -i; x <= i; x++ {
			if w.loadChunkIfNotLoaded(x*ChunkSize+chunkX, 0, z*ChunkSize+chunkZ) {
				for y := 1; y < WorldHeight/ChunkSize; y++ {
					w.LoadChunk(x*ChunkSize+chunkX, y*ChunkSize, z*ChunkSize+chunkZ)
				}
				return
			}
		}
		x := i
		for z := -i + 1; z < i; z++ {
			if w.loadChunkIfNotLoaded(x*ChunkSize+chunkX, 0, z*ChunkSize) {
				for y := 1; y < WorldHeight/ChunkSize; y++ {
					w.LoadChunk(x*ChunkSize+chunkX, y*ChunkSize, z*ChunkSize+chunkZ)
				}
				return
			}
		}
		z = i
		for x := i; x >= -i; x-- {
			if w.loadChunkIfNotLoaded(x*ChunkSize+chunkX, 0, z*ChunkSize) {
				for y := 1; y < WorldHeight/ChunkSize; y++ {
					w.LoadChunk(x*ChunkSize+chunkX, y*ChunkSize, z*ChunkSize+chunkZ)
				}
				return
			}
		}
		x = -i
		for z := i - 1; z > -i; z-- {
			if w.loadChunkIfNotLoaded(x*ChunkSize+chunkX, 0, z*ChunkSize+chunkZ) {
				for y := 1; y < WorldHeight/ChunkSize; y++ {
					w.LoadChunk(x*ChunkSize+chunkX, y*ChunkSize, z*ChunkSize+chunkZ)
				}
				return
			}
		}
	}
}

//LoadChunks2 load chunks around the player, for now, only load one chunk per frame
func (w *World) LoadChunks2(playerPos mgl32.Vec3) {
	xPlayer := int(playerPos.X())
	zPlayer := int(playerPos.Z())
	chunkX := getChunk(xPlayer)
	chunkZ := getChunk(zPlayer)
	for i := 0; i < chunkMaxLoadDistance; i++ {
		p := Point{i * ChunkSize, 0, i * ChunkSize}
		if p.DistanceTo(playerPos) > loadChunkDistance {
			return
		}
		z := -i
		for x := -i; x <= i; x++ {
			if w.loadChunkIfNotLoaded(x*ChunkSize+chunkX, 0, z*ChunkSize+chunkZ) {
				for y := 1; y < WorldHeight/ChunkSize; y++ {
					w.LoadChunk(x*ChunkSize+chunkX, y*ChunkSize, z*ChunkSize+chunkZ)
				}
			}
		}
		x := i
		for z := -i + 1; z < i; z++ {
			if w.loadChunkIfNotLoaded(x*ChunkSize+chunkX, 0, z*ChunkSize) {
				for y := 1; y < WorldHeight/ChunkSize; y++ {
					w.LoadChunk(x*ChunkSize+chunkX, y*ChunkSize, z*ChunkSize+chunkZ)
				}
			}
		}
		z = i
		for x := i; x >= -i; x-- {
			if w.loadChunkIfNotLoaded(x*ChunkSize+chunkX, 0, z*ChunkSize) {
				for y := 1; y < WorldHeight/ChunkSize; y++ {
					w.LoadChunk(x*ChunkSize+chunkX, y*ChunkSize, z*ChunkSize+chunkZ)
				}
			}
		}
		x = -i
		for z := i - 1; z > -i; z-- {
			if w.loadChunkIfNotLoaded(x*ChunkSize+chunkX, 0, z*ChunkSize+chunkZ) {
				for y := 1; y < WorldHeight/ChunkSize; y++ {
					w.LoadChunk(x*ChunkSize+chunkX, y*ChunkSize, z*ChunkSize+chunkZ)
				}
			}
		}
	}
}
