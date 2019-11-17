package world

import (
	"fmt"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/piochelepiotr/minecraftGo/entities"
	"github.com/piochelepiotr/minecraftGo/geometry"
	"github.com/piochelepiotr/minecraftGo/loader"
	"github.com/piochelepiotr/minecraftGo/models"
	"github.com/piochelepiotr/minecraftGo/render"
	"github.com/piochelepiotr/minecraftGo/textures"
	"github.com/piochelepiotr/minecraftGo/toolbox"
	"math"
	"time"
)


const (
	deleteChunkDistance  float32 = 120
	loadChunkDistance    float32 = 100
	maxWallJump float32 = 0.4
	backwardJump float32 = 0.4
	playerWidth float32 = 0.4
	playerVWidth = playerWidth/1.5
	playerHeight float32 = 1.8
)

var cosAngle = float32(math.Cos(float64(render.Fov/2 + mgl32.DegToRad(10))))
var alwaysRenderDistance = float32(20)//float32(ChunkSize)/ float32( 2 * math.Tan(float64(render.Fov/2)))

// World contains all the blocks of the world in chunks that load around the player
type World struct {
	chunks       map[geometry.Point]*Chunk
	modelTexture textures.ModelTexture
	generator *Generator
	ChunkLoadDecisions chan geometry.Point
}

func getChunk(x int) int {
	return int(math.Floor(float64(x)/float64(ChunkSize))) * ChunkSize
}

// CreateWorld initiate the world
func CreateWorld(generator *Generator) *World {
	modelTexture := loader.LoadModelTexture("textures/textures2.png", 16)
	chunks := make(map[geometry.Point]*Chunk)
	fmt.Println(alwaysRenderDistance)
	return &World{
		chunks:       chunks,
		modelTexture: modelTexture,
		generator: generator,
		ChunkLoadDecisions: make(chan geometry.Point, 200),
	}
}

func isVisible(cameraPosition, coneVector, pos mgl32.Vec3) bool {
	if cameraPosition.Sub(pos).Len() < alwaysRenderDistance {
		return true
	}
	toPoint := pos.Sub(cameraPosition).Normalize()
	dotProduct := toPoint.Dot(coneVector)
	return dotProduct > cosAngle
}
func isChunkVisible(cameraPosition, coneVector, corner mgl32.Vec3) bool {
	points := []mgl32.Vec3{
		{0, 0, 0},
		{float32(ChunkSize), 0, 0},
		{0, float32(ChunkSize), 0},
		{0, 0, float32(ChunkSize)},
		{float32(ChunkSize), float32(ChunkSize), 0},
		{float32(ChunkSize), 0, float32(ChunkSize)},
		{0, float32(ChunkSize), float32(ChunkSize)},
		{float32(ChunkSize), float32(ChunkSize), float32(ChunkSize)},
	}
	for _, p := range points {
		if isVisible(cameraPosition, coneVector, corner.Add(p)) {
			return true
		}
	}
	return false
}

//GetChunks returns all chunks that are going to be rendered
func (w *World) GetChunks(camera *entities.Camera) []entities.Entity {
	chunks := make([]entities.Entity, 0)
	for _, chunk := range w.chunks {
		model := chunk.Model
		if model.VertexCount == 0 {
			continue
		}
		coneVector := geometry.ComputeCameraRay(camera.Rotation).Normalize()
		p := mgl32.Vec3{
				float32(chunk.Start.X),
				float32(chunk.Start.Y),
				float32(chunk.Start.Z),
		}
		if !isChunkVisible(camera.Position, coneVector, p) {
			continue
		}
		chunkEntity := entities.Entity{
			TexturedModel: models.TexturedModel{
				RawModel:     chunk.Model,
				ModelTexture: w.modelTexture,
			},
			Position: p,
		}
		transparentChunkEntity := entities.Entity{
			TexturedModel: models.TexturedModel{
				RawModel:     chunk.TransparentModel,
				ModelTexture: w.modelTexture,
				Transparent: true,
			},
			Position: mgl32.Vec3{
				float32(chunk.Start.X),
				float32(chunk.Start.Y),
				float32(chunk.Start.Z),
			},
		}
		if chunkEntity.TexturedModel.RawModel.VertexCount > 0 {
			chunks = append(chunks, chunkEntity)
		}
		if transparentChunkEntity.TexturedModel.RawModel.VertexCount > 0 {
			chunks = append(chunks, transparentChunkEntity)
		}
	}
	return chunks
}

// LoadChunk loads a chunk into the world so that it is rendered
func (w *World) LoadChunk(x, y, z int) {
	p := geometry.Point{
		X: x,
		Y: y,
		Z: z,
	}
	w.chunks[p] = CreateChunk(p, w.generator)
}

func (w *World) chunkIsLoaded(x, y, z int) bool {
	_, ok := w.chunks[geometry.Point{x, y, z}]
	return ok
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
	p := geometry.Point{
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
	p := geometry.Point{
		X: chunkX,
		Y: chunkY,
		Z: chunkZ,
	}
	if chunk, ok := w.chunks[p]; ok {
		chunk.SetBlock(x-chunkX, y-chunkY, z-chunkZ, b)
	} else {
		fmt.Println("ERROR when setting block in chunk ", p)
	}
}


// even if the point is a bit inside a wall, this is going to return
func (w *World) PlaceInFrontWithJumps(edges []mgl32.Vec3, dir mgl32.Vec3) float32 {
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
func (w *World) PlaceInFrontWithJumpsOnePoint(p mgl32.Vec3, dir mgl32.Vec3) float32 {
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
func (w *World) PlaceInFront(p mgl32.Vec3, dir mgl32.Vec3) (float32, geometry.Point, geometry.Point) {
	xray := geometry.NewXray(p, dir)
	for i := 0; i < 10; i++ {
		//fmt.Println("POS: ", x, ";", y, ";", z)
		if w.GetBlock(xray.P.X, xray.P.Y, xray.P.Z) != Air {
			return xray.Distance, xray.P, xray.Previous
		}
		xray.GoToNextBlock()
	}
	return xray.Distance, geometry.Point{}, geometry.Point{}
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
	return w.PlaceInFrontWithJumps(returnPlayerVerticalEdges(player), mgl32.Vec3{0, -1, 0}) == 0
}


func returnPlayerEdges(player *entities.Player) []mgl32.Vec3 {
	forward := player.FacingDir(playerWidth)
	backward := forward.Mul(-1)
	side1 := player.SideFacingDir(playerWidth)
	side2 := side1.Mul(-1)
	up := mgl32.Vec3{0, playerHeight, 0}
	edges := make([]mgl32.Vec3,0, 8)
	edges = append(edges, player.Entity.Position.Add(forward).Add(side1))
	edges = append(edges, player.Entity.Position.Add(forward).Add(side2))
	edges = append(edges, player.Entity.Position.Add(backward).Add(side1))
	edges = append(edges, player.Entity.Position.Add(backward).Add(side2))
	for i := 0; i < 4; i++ {
		edges = append(edges, edges[i].Add(up))
	}
	return edges
}

func returnPlayerVerticalEdges(player *entities.Player) []mgl32.Vec3 {
	forward := player.FacingDir(playerVWidth)
	backward := forward.Mul(-1)
	side1 := player.SideFacingDir(playerVWidth)
	side2 := side1.Mul(-1)
	up := mgl32.Vec3{0, playerHeight, 0}
	edges := make([]mgl32.Vec3,0, 8)
	edges = append(edges, player.Entity.Position.Add(forward).Add(side1))
	edges = append(edges, player.Entity.Position.Add(forward).Add(side2))
	edges = append(edges, player.Entity.Position.Add(backward).Add(side1))
	edges = append(edges, player.Entity.Position.Add(backward).Add(side2))
	for i := 0; i < 4; i++ {
		edges = append(edges, edges[i].Add(up))
	}
	return edges
}


func (w *World) PlacePlayerOnGround(player *entities.Player) {
	p := player.Entity.Position
	dir := float32(-1)
	for {
		edges := returnPlayerVerticalEdges(player)
		place := w.PlaceInFrontWithJumps(edges, mgl32.Vec3{0, dir, 0})
		if place == 0 {
			return
		}
		player.Entity.Position = mgl32.Vec3{p.X(), float32(math.Floor(float64(dir*place+player.Entity.Position.Y()))), p.Z()}
	}
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
			edges := returnPlayerVerticalEdges(player)
			if y > 0 {
				dir = 1
				edges = edges[4:]
			} else {
				edges = edges[:4]
			}
			place := w.PlaceInFrontWithJumps(edges, mgl32.Vec3{0, dir, 0})
			// fmt.Printf("Speed y is %f\n", player.Speed.Y())
			if toolbox.Abs(y) > place {
				player.Speed = mgl32.Vec3{player.Speed.X(), 0, player.Speed.Z()}
				newY := dir*place+player.Entity.Position.Y()
				if dir == -1 {
					newY = float32(math.Floor(float64(newY)))
				}
				player.Entity.Position = mgl32.Vec3{p.X(), newY, p.Z()}
			} else {
				player.Entity.Position = player.Entity.Position.Add(mgl32.Vec3{0, y, 0})
			}
			// fmt.Printf("Pos y is %f\n", player.Entity.Position.Y())
		}

		hSpeed := mgl32.Vec3{player.Speed.X(), 0, player.Speed.Z()}

		// horizontal movement second
		if hSpeed.Len() > 0 {
			move := hSpeed.Mul(dt)
			// fmt.Printf("Move is %f\n", move.Len())
			//forward := mgl32.Vec3{0, -1, 0}
			forward := hSpeed
			// first, go as far as we can in the forward direction
			place := w.PlaceInFrontWithJumps(returnPlayerEdges(player), forward)
			firstMove := truncateMovement(move, place)
			// fmt.Printf("first move is %f\n", firstMove.Len())
			player.Entity.Position = player.Entity.Position.Add(firstMove)
			restMove := move.Sub(firstMove)

			if restMove.Len() > 0 {
				if restMove.X() != 0 {
					placeX := w.PlaceInFrontWithJumps(returnPlayerEdges(player), mgl32.Vec3{hSpeed.X(), 0, 0})
					if placeX > toolbox.Abs(restMove.X()) {
						player.Entity.Position = player.Entity.Position.Add(mgl32.Vec3{restMove.X(), 0, 0})
					} else {
						player.Speed = mgl32.Vec3{0, player.Speed.Y(), player.Speed.Z()}
						if placeX > 0 {
							player.Entity.Position = player.Entity.Position.Add(mgl32.Vec3{(restMove.X()/toolbox.Abs(restMove.X()))*placeX, 0, 0})
						}
					}
				}
				if restMove.Z() != 0 {
					placeZ := w.PlaceInFrontWithJumps(returnPlayerEdges(player), mgl32.Vec3{0, 0, hSpeed.Z()})
					if placeZ > toolbox.Abs(restMove.Z()) {
						player.Entity.Position = player.Entity.Position.Add(mgl32.Vec3{0, 0, restMove.Z()})
					} else {
						player.Speed = mgl32.Vec3{player.Speed.X(), player.Speed.Y(), 0}
						if placeZ > 0 {
							player.Entity.Position = player.Entity.Position.Add(mgl32.Vec3{0, 0, (restMove.Z()/toolbox.Abs(restMove.Z()))*placeZ})
						}
					}
				}
			}
		}
		if !(forward || backward) {
			player.Speed = entities.Friction(player.Speed, dt)
		}
		player.Speed = entities.Gravity(player.Speed, dt, touchGround)
	}
	player.LastMove = now
}

// ClickOnBlock removes or adds a block
func (w *World) ClickOnBlock(camera *entities.Camera, placeBlock bool) {
	xray := geometry.ComputeCameraRay(camera.Rotation)
	p := camera.Position
	_, block, previous := w.PlaceInFront(p, xray)
	if placeBlock {
		if !block.Equal(previous) {
			w.SetBlock(previous.X, previous.Y, previous.Z, Tree)
		}
	} else {
		w.SetBlock(block.X, block.Y, block.Z, Air)
	}
}

func (w *World) deleteChunks(playerPos mgl32.Vec3) {
	for p, chunk := range w.chunks {
		if chunk.Start.DistanceTo(playerPos) > deleteChunkDistance {
			delete(w.chunks, p)
		}
	}
}

//LoadChunks load chunks around the player
func (w *World) LoadChunks(playerPos mgl32.Vec3, delay bool) {
	w.deleteChunks(playerPos)
	xPlayer := int(playerPos.X())
	zPlayer := int(playerPos.Z())
	chunkX := getChunk(xPlayer)
	chunkZ := getChunk(zPlayer)
	for x := getChunk(chunkX - int(loadChunkDistance)); x <= getChunk(chunkX + int(loadChunkDistance)); x += ChunkSize {
		for z := getChunk(chunkZ - int(loadChunkDistance)); z <= getChunk(chunkZ + int(loadChunkDistance)); z += ChunkSize {
			p := geometry.Point{x, 0, z}
			if p.DistanceTo(playerPos) > loadChunkDistance {
				continue
			}
			if w.chunkIsLoaded(x, 0, z) {
				continue
			}
			for y := 0; y < WorldHeight/ChunkSize; y++ {
				if delay {
					w.ChunkLoadDecisions <- geometry.Point{x, y*ChunkSize, z}
				} else {
					w.LoadChunk(x, y*ChunkSize, z)
				}
			}
		}
	}
}

func (w *World) AddChunk(chunk *Chunk) {
	w.chunks[chunk.Start] = chunk
}

// Close saves the world when closing the game
func (w *World) Close() {
	close(w.ChunkLoadDecisions)
}
