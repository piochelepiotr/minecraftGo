package world

import (
	"log"
	"math"
	"sort"
	"time"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/piochelepiotr/minecraftGo/entities"
	"github.com/piochelepiotr/minecraftGo/game_engine/loader"
	"github.com/piochelepiotr/minecraftGo/game_engine/models"
	"github.com/piochelepiotr/minecraftGo/game_engine/render"
	"github.com/piochelepiotr/minecraftGo/geometry"
	"github.com/piochelepiotr/minecraftGo/textures"
	"github.com/piochelepiotr/minecraftGo/toolbox"
	"github.com/piochelepiotr/minecraftGo/world/block"
	"github.com/piochelepiotr/minecraftGo/worldcontent"
)

const (
	playerWidth  float32 = 0.4
	playerVWidth         = playerWidth / 1.5
	playerHeight float32 = 1.8
)

// todo: fov isn't the same vertically and horizontally, so something is wrong with that. It depends on
// aspectRatio
var alwaysRenderDistance = float32(20) //float32(ChunkSize)/ float32( 2 * math.Tan(float64(render.Fov/2)))

// World contains all the blocks of the world in chunks that load around the player
type World struct {
	chunks map[geometry.Point]*Chunk
	// loading is used to not load two times the same chunk
	loading      map[geometry.Point]struct{}
	modelTexture textures.ModelTexture
	chunksToLoad chan geometry.Point
	cosFovAngle  float32
	world        *worldcontent.InMemoryWorld
}

// NewWorld initiate the world
func NewWorld(world *worldcontent.InMemoryWorld, aspectRatio float32) *World {
	modelTexture := loader.LoadModelTexture("textures/textures2.png", uint32(block.NumberRowsTextures))
	chunks := make(map[geometry.Point]*Chunk)
	w := &World{
		chunks:       chunks,
		loading:      make(map[geometry.Point]struct{}),
		modelTexture: modelTexture,
		// really big here to avoid dead locks (main goroutines writes to this channel
		// while loader reads from this one and writes a a channel the main goroutine reads from
		chunksToLoad: make(chan geometry.Point, 50000),
		world:        world,
	}
	w.Resize(aspectRatio)
	return w
}

func (w *World) ChunksToLoad() <-chan geometry.Point {
	return w.chunksToLoad
}

func (w *World) Resize(aspectRatio float32) {
	fovY := render.Fov
	fovX := fovY * aspectRatio
	// fov := float32(math.Max(float64(fovX), float64(fovY)))
	fov := float32(math.Sqrt(math.Pow(float64(fovX), 2) + math.Pow(float64(fovY), 2)))
	// log.Println("fov", mgl32.RadToDeg(fov))
	w.cosFovAngle = float32(math.Cos(float64(fov / 2)))
}

func (w *World) isVisible(cameraPosition, coneVector, pos mgl32.Vec3) bool {
	if cameraPosition.Sub(pos).Len() < alwaysRenderDistance {
		return true
	}
	toPoint := pos.Sub(cameraPosition).Normalize()
	dotProduct := toPoint.Dot(coneVector)
	return dotProduct > w.cosFovAngle
}
func (w *World) isChunkVisible(cameraPosition, coneVector, corner mgl32.Vec3) bool {
	points := []mgl32.Vec3{
		{0, 0, 0},
		{float32(worldcontent.ChunkSize), 0, 0},
		{0, float32(worldcontent.ChunkSize), 0},
		{0, 0, float32(worldcontent.ChunkSize)},
		{float32(worldcontent.ChunkSize), float32(worldcontent.ChunkSize), 0},
		{float32(worldcontent.ChunkSize), 0, float32(worldcontent.ChunkSize)},
		{0, float32(worldcontent.ChunkSize), float32(worldcontent.ChunkSize)},
		{float32(worldcontent.ChunkSize), float32(worldcontent.ChunkSize), float32(worldcontent.ChunkSize)},
	}
	for _, p := range points {
		if w.isVisible(cameraPosition, coneVector, corner.Add(p)) {
			return true
		}
	}
	return false
}

//GetChunks returns all chunks that are going to be rendered
// no need to render chunks behind the player
func (w *World) GetChunks(camera *entities.Camera) []entities.Entity {
	chunks := make([]entities.Entity, 0)
	for _, chunk := range w.chunks {
		model := chunk.Model
		if model.VertexCount == 0 {
			continue
		}
		coneVector := geometry.ComputeCameraRay(camera.Rotation).Normalize()
		p := mgl32.Vec3{
			float32(chunk.start.X),
			float32(chunk.start.Y),
			float32(chunk.start.Z),
		}
		if !w.isChunkVisible(camera.Position, coneVector, p) {
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
				Transparent:  true,
			},
			Position: mgl32.Vec3{
				float32(chunk.start.X),
				float32(chunk.start.Y),
				float32(chunk.start.Z),
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

func (w *World) chunkIsLoaded(x, y, z int) bool {
	_, ok := w.chunks[geometry.Point{x, y, z}]
	return ok
}

func (w *World) chunkIsLoading(x, y, z int) bool {
	_, ok := w.loading[geometry.Point{x, y, z}]
	return ok
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
	return w.world.PlaceInFrontWithJumps(returnPlayerVerticalEdges(player), mgl32.Vec3{0, -1, 0}) == 0
}

func returnPlayerEdges(player *entities.Player) []mgl32.Vec3 {
	forward := player.FacingDir(playerWidth)
	backward := forward.Mul(-1)
	side1 := player.SideFacingDir(playerWidth)
	side2 := side1.Mul(-1)
	up := mgl32.Vec3{0, playerHeight, 0}
	edges := make([]mgl32.Vec3, 0, 8)
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
	edges := make([]mgl32.Vec3, 0, 8)
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
		place := w.world.PlaceInFrontWithJumps(edges, mgl32.Vec3{0, dir, 0})
		if place == 0 {
			return
		}
		player.Entity.Position = mgl32.Vec3{p.X(), float32(math.Floor(float64(dir*place + player.Entity.Position.Y()))), p.Z()}
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
			place := w.world.PlaceInFrontWithJumps(edges, mgl32.Vec3{0, dir, 0})
			// fmt.Printf("Speed y is %f\n", player.Speed.Y())
			if toolbox.Abs(y) > place {
				player.Speed = mgl32.Vec3{player.Speed.X(), 0, player.Speed.Z()}
				newY := dir*place + player.Entity.Position.Y()
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
			place := w.world.PlaceInFrontWithJumps(returnPlayerEdges(player), forward)
			firstMove := truncateMovement(move, place)
			// fmt.Printf("first move is %f\n", firstMove.Len())
			player.Entity.Position = player.Entity.Position.Add(firstMove)
			restMove := move.Sub(firstMove)

			if restMove.Len() > 0 {
				if restMove.X() != 0 {
					placeX := w.world.PlaceInFrontWithJumps(returnPlayerEdges(player), mgl32.Vec3{hSpeed.X(), 0, 0})
					if placeX > toolbox.Abs(restMove.X()) {
						player.Entity.Position = player.Entity.Position.Add(mgl32.Vec3{restMove.X(), 0, 0})
					} else {
						player.Speed = mgl32.Vec3{0, player.Speed.Y(), player.Speed.Z()}
						if placeX > 0 {
							player.Entity.Position = player.Entity.Position.Add(mgl32.Vec3{(restMove.X() / toolbox.Abs(restMove.X())) * placeX, 0, 0})
						}
					}
				}
				if restMove.Z() != 0 {
					placeZ := w.world.PlaceInFrontWithJumps(returnPlayerEdges(player), mgl32.Vec3{0, 0, hSpeed.Z()})
					if placeZ > toolbox.Abs(restMove.Z()) {
						player.Entity.Position = player.Entity.Position.Add(mgl32.Vec3{0, 0, restMove.Z()})
					} else {
						player.Speed = mgl32.Vec3{player.Speed.X(), player.Speed.Y(), 0}
						if placeZ > 0 {
							player.Entity.Position = player.Entity.Position.Add(mgl32.Vec3{0, 0, (restMove.Z() / toolbox.Abs(restMove.Z())) * placeZ})
						}
					}
				}
			}
		}
		if !(forward || backward) {
			player.Speed = entities.Friction(player.Speed, dt)
		}
		if !player.InFlight {
			player.Speed = entities.Gravity(player.Speed, dt, touchGround)
		}
	}
	player.LastMove = now
}

func (w *World) rebuild(x, y, z int) {
	if chunk, ok := w.chunks[geometry.Point{x, y, z}]; ok {
		chunk.buildFaces()
		chunk.Load()
	}
}

// setBlock sets a block and updates UI if necessary
func (w *World) setBlock(x, y, z int, b block.Block) {
	log.Println("setting block", x, y, z)
	updated := w.world.SetBlock(x, y, z, b)
	if !updated {
		return
	}
	chunkX := worldcontent.ChunkStart(x)
	chunkY := worldcontent.ChunkStart(y)
	chunkZ := worldcontent.ChunkStart(z)
	w.rebuild(chunkX, chunkY, chunkZ)
	if chunkX2 := worldcontent.ChunkStart(x - 1); chunkX2 != chunkX {
		w.rebuild(chunkX2, chunkY, chunkZ)
	}
	if chunkX2 := worldcontent.ChunkStart(x + 1); chunkX2 != chunkX {
		w.rebuild(chunkX2, chunkY, chunkZ)
	}
	if chunkY2 := worldcontent.ChunkStart(y - 1); chunkY2 != chunkX {
		w.rebuild(chunkX, chunkY2, chunkZ)
	}
	if chunkY2 := worldcontent.ChunkStart(y + 1); chunkY2 != chunkX {
		w.rebuild(chunkX, chunkY2, chunkZ)
	}
	if chunkZ2 := worldcontent.ChunkStart(z - 1); chunkZ2 != chunkX {
		w.rebuild(chunkX, chunkY, chunkZ2)
	}
	if chunkZ2 := worldcontent.ChunkStart(z + 1); chunkZ2 != chunkX {
		w.rebuild(chunkX, chunkY, chunkZ2)
	}
}

// ClickOnBlock removes or adds a block
func (w *World) ClickOnBlock(camera *entities.Camera, placeBlock bool, b block.Block) {
	xray := geometry.ComputeCameraRay(camera.Rotation)
	p := camera.Position
	_, pos, previous := w.world.PlaceInFront(p, xray)
	if placeBlock {
		if !pos.Equal(previous) {
			w.setBlock(previous.X, previous.Y, previous.Z, b)
		}
	} else {
		w.setBlock(pos.X, pos.Y, pos.Z, block.Air)
	}
}

func (w *World) deleteChunks(playerPos mgl32.Vec3) {
	for p, chunk := range w.chunks {
		if chunk.start.DistanceTo(playerPos) > worldcontent.UIDeleteChunkDistance {
			delete(w.chunks, p)
		}
	}
}

type chunkToLoad struct {
	distance float32
	p        geometry.Point
}

type chunksToLoad []chunkToLoad

func (c chunksToLoad) Len() int {
	return len(c)
}

func (c chunksToLoad) Less(i, j int) bool {
	return c[i].distance < c[j].distance
}

func (c chunksToLoad) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

//LoadChunks load chunks around the player
func (w *World) LoadChunks(playerPos mgl32.Vec3) {
	w.deleteChunks(playerPos)
	w.world.ExpireChunks(playerPos)
	xPlayer := int(playerPos.X())
	zPlayer := int(playerPos.Z())
	chunkX := worldcontent.ChunkStart(xPlayer)
	chunkZ := worldcontent.ChunkStart(zPlayer)
	toLoad := make(chunksToLoad, 0)
	for x := worldcontent.ChunkStart(chunkX - int(worldcontent.UILoadChunkDistance)); x <= worldcontent.ChunkStart(chunkX+int(worldcontent.UILoadChunkDistance)); x += worldcontent.ChunkSize {
		for z := worldcontent.ChunkStart(chunkZ - int(worldcontent.UILoadChunkDistance)); z <= worldcontent.ChunkStart(chunkZ+int(worldcontent.UILoadChunkDistance)); z += worldcontent.ChunkSize {
			p := geometry.Point{X: x, Y: 0, Z: z}
			d := p.DistanceTo(playerPos)
			if d > worldcontent.UILoadChunkDistance {
				continue
			}
			for y := 0; y < worldcontent.WorldHeight; y += worldcontent.ChunkSize {
				if w.chunkIsLoaded(x, y, z) || w.chunkIsLoading(x, y, z) {
					continue
				}
				toLoad = append(toLoad, chunkToLoad{distance: d, p: geometry.Point{x, y, z}})
				// select {
				// 	case w.chunksToLoad <- geometry.Point{x, y*w.world.ChunkSize(), z}:
				// 	default:
				// 		log.Println("couldn't load chunk. Retrying soon")
				// }
			}
		}
	}
	sort.Sort(toLoad)
	for _, chunk := range toLoad {
		w.loading[chunk.p] = struct{}{}
		w.chunksToLoad <- chunk.p
	}
}

func (w *World) AddChunk(chunk *Chunk) {
	chunk.Load()
	w.chunks[chunk.start] = chunk
	delete(w.loading, chunk.start)
}

// Close saves the world when closing the game
func (w *World) Close() {
	close(w.chunksToLoad)
	log.Println("Closed world")
}
