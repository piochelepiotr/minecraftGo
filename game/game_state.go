package game

import (
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/piochelepiotr/minecraftGo/entities"
	"github.com/piochelepiotr/minecraftGo/loader"
	"github.com/piochelepiotr/minecraftGo/models"
	"github.com/piochelepiotr/minecraftGo/render"
	"github.com/piochelepiotr/minecraftGo/state"
	pworld "github.com/piochelepiotr/minecraftGo/world"
)

type keyPressed struct {
	wPressed     bool
	aPressed     bool
	dPressed     bool
	sPressed     bool
	spacePressed bool
}

// GamingState is the 3D state
type GamingState struct {
	world       *pworld.World
	player      *entities.Player
	camera      *entities.Camera
	light       *entities.Light
	keyPressed  keyPressed
	chunkLoader *pworld.ChunkLoader
	settings *settings
	doneWriter <-chan struct{}
	display     *render.DisplayManager
	changeState chan<- state.ID
}
// NewGamingState loads a new world
func NewGamingState(display *render.DisplayManager, changeState chan<- state.ID) *GamingState{
	generator := pworld.NewGenerator()
	chunkLoader := pworld.NewChunkLoader(generator)
	world := pworld.CreateWorld(generator)
	doneWriter := pworld.NewChunkWriter(world.OutChunksToWrite())
	chunkLoader.Run(world.ChunkLoadDecisions)
	camera := entities.CreateCamera(-50, 30, -50, -0.2, 1.8)
	camera.Rotation = mgl32.Vec3{0, 0, 0}

	model := loader.LoadObjModel("objects/steve.obj")
	t := loader.LoadModelTexture("textures/skin.png", 1)
	texturedModel := models.TexturedModel{
		ModelTexture: t,
		RawModel:     model,
	}

	light := &entities.Light{
		Position: mgl32.Vec3{5, 5, 5},
		Colour:   mgl32.Vec3{1, 1, 1},
	}

	entity := entities.Entity{
		Position:      mgl32.Vec3{0, float32(pworld.WorldHeight + 20), 0},
		Rotation:      mgl32.Vec3{0, 0, 0},
		TexturedModel: texturedModel,
	}
	player := &entities.Player{
		Entity: entity,
	}
	world.LoadChunks(player.Entity.Position, false)
	world.PlacePlayerOnGround(player)

	state := &GamingState{
		world:       world,
		player:      player,
		camera:      camera,
		light:       light,
		chunkLoader: chunkLoader,
		settings: defaultSettings(),
		doneWriter: doneWriter,
		display: display,
		changeState: changeState,
	}
	return state
}

func (s *GamingState) Close() {
	s.world.Close()
	<-s.doneWriter
}

func (s *GamingState) clickCallback(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
	if action == glfw.Press {
		if button == glfw.MouseButtonRight {
			s.world.ClickOnBlock(s.camera, true)
		} else if button == glfw.MouseButtonLeft {
			s.world.ClickOnBlock(s.camera, false)
		}
	}
}

func (s *GamingState) mouseMoveCallback(w *glfw.Window, xpos float64, ypos float64) {
	x, y := s.display.GLPos(xpos, ypos)
	s.player.Entity.Rotation = mgl32.Vec3{0, -x*s.settings.cameraSensitivity, 0}
	s.camera.Rotation = mgl32.Vec3{y*s.settings.cameraSensitivity, s.camera.Rotation.Y(), s.camera.Rotation.Z()}
}

func (s *GamingState) keyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if action == glfw.Press {
		if key == glfw.KeyW {
			s.keyPressed.wPressed = true
		} else if key == glfw.KeyA {
			s.keyPressed.aPressed = true
		} else if key == glfw.KeyD {
			s.keyPressed.dPressed = true
		} else if key == glfw.KeyS {
			s.keyPressed.sPressed = true
		} else if key == glfw.KeySpace {
			s.keyPressed.spacePressed = true
		} else if key == glfw.KeyEscape {
			s.changeState <- state.GameMenu
		}
	} else if action == glfw.Release {
		if key == glfw.KeyW {
			s.keyPressed.wPressed = false
		} else if key == glfw.KeyA {
			s.keyPressed.aPressed = false
		} else if key == glfw.KeyD {
			s.keyPressed.dPressed = false
		} else if key == glfw.KeyS {
			s.keyPressed.sPressed = false
		} else if key == glfw.KeySpace {
			s.keyPressed.spacePressed = false
		}
	}
}
// Render renders all objects on the screen
func (s *GamingState) Render(renderer *render.MasterRenderer) {
	// why here?
	s.camera.LockOnPlayer(s.player)
	// r.ProcessEntity(player.Entity)
	renderer.ProcessEntities(s.world.GetChunks(s.camera))
	renderer.Render(s.light, s.camera)
}
// NextFrame makes time pass to move to the next frame of the game
func (s *GamingState) NextFrame() {
	select {
	case chunk := <-s.chunkLoader.LoadedChunk:
		s.world.AddChunk(chunk)
	default:
	}
	forward := s.keyPressed.wPressed
	backward := s.keyPressed.sPressed
	right := s.keyPressed.dPressed
	left := s.keyPressed.aPressed
	jump := s.keyPressed.spacePressed
	touchGround := s.world.TouchesGround(s.player)
	s.player.Move(forward, backward, jump, touchGround, right, left)
	s.world.MovePlayer(s.player, forward, backward, jump, touchGround)
}
// Update is called every second
func (s *GamingState) Update() {
	s.world.LoadChunks(s.player.Entity.Position, true)
}

// pause is used when menu is opened
func (s *GamingState) pause() {
	s.keyPressed = keyPressed{}
}