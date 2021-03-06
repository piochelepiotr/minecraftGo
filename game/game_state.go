package game

import (
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/piochelepiotr/minecraftGo/entities"
	"github.com/piochelepiotr/minecraftGo/game_engine/guis"
	"github.com/piochelepiotr/minecraftGo/game_engine/loader"
	"github.com/piochelepiotr/minecraftGo/game_engine/models"
	"github.com/piochelepiotr/minecraftGo/game_engine/render"
	"github.com/piochelepiotr/minecraftGo/state"
	"github.com/piochelepiotr/minecraftGo/ux"
	pworld "github.com/piochelepiotr/minecraftGo/world"
	"github.com/piochelepiotr/minecraftGo/world/block"
	"github.com/piochelepiotr/minecraftGo/worldcontent"
	"log"
	"time"
)

type keyPressed struct {
	previousSpacePressed time.Time
	wPressed     bool
	aPressed     bool
	dPressed     bool
	sPressed     bool
	spacePressed bool
	shiftPressed bool
}

// GamingState is the 3D state
type GamingState struct {
	cubeOutline models.OutlineModel
	worldContent *worldcontent.InMemoryWorld
	world       *pworld.World
	player      *entities.Player
	camera      *entities.Camera
	cursor      guis.GuiTexture
	keyPressed  keyPressed
	chunkLoader *pworld.ChunkLoader
	settings *settings
	doneWriter <-chan struct{}
	display     *render.DisplayManager
	changeState chan<- state.Switch
	bottomBar *ux.BottomBar
	scroll float64
	inventory *pworld.Inventory
}
// NewGamingState loads a new world
func NewGamingState(worldName string, display *render.DisplayManager, changeState chan<- state.Switch, loader *loader.Loader, structEditing bool) *GamingState{
	worldConfig, err := worldcontent.LoadWorld(worldName)
	if err != nil {
		log.Fatalf("Unable to load world %s. Err: %v", worldName, err)
	}
	inventory := pworld.NewInventory()
	wContent := worldcontent.NewInMemoryWorld(worldConfig, structEditing)
	world := pworld.NewWorld(wContent, display.AspectRatio(), loader)
	chunkLoader := pworld.NewChunkLoader(wContent, world.ChunksToLoad())
	doneWriter := worldcontent.NewChunkWriter(worldConfig, wContent.OutChunksToWrite())
	camera := entities.CreateCamera(-50, 30, -50, -0.2, 1.8)
	camera.Rotation = mgl32.Vec3{0, 0, 0}

	model := loader.LoadObjModel("objects/steve.obj")
	textureID := loader.LoadModelTexture("textures/skin.png")
	texturedModel := models.TexturedModel{
		TextureID: textureID,
		RawModel:  model,
	}

	entity := entities.Entity{
		Position:      mgl32.Vec3{worldConfig.Player.PosX, worldConfig.Player.PosY, worldConfig.Player.PosZ},
		Rotation:      mgl32.Vec3{0, 0, 0},
		TexturedModel: texturedModel,
	}
	player := &entities.Player{
		Entity: entity,
	}
	if player.Entity.Position.Y() == -1 {
		player.Entity.Position = mgl32.Vec3{player.Entity.Position.X(), float32(worldcontent.WorldHeight + 20), player.Entity.Position.Z()}
		world.PlacePlayerOnGround(player)
	}

	state := &GamingState{
		worldContent: wContent,
		cursor:      loader.LoadGuiTexture("textures/cursor.png", mgl32.Vec2{0, 0}, mgl32.Vec2{0.02, 0.03}),
		world:       world,
		player:      player,
		camera:      camera,
		chunkLoader: chunkLoader,
		settings: defaultSettings(),
		doneWriter: doneWriter,
		display: display,
		changeState: changeState,
		bottomBar: ux.NewBottomBar(display.AspectRatio(), loader, inventory),
		cubeOutline: ux.NewCubeOutline(loader),
		inventory: inventory,
	}
	world.LoadChunks(player.Entity.Position)
	// state.loadChunks(player.Entity.Position)
	return state
}

func (s *GamingState) Close() {
	s.world.Close()
	s.worldContent.Close()
	<-s.doneWriter
	pos := s.player.Entity.Position
	cfg := s.worldContent.Config()
	cfg.Player.PosX = pos.X()
	cfg.Player.PosY = pos.Y()
	cfg.Player.PosZ = pos.Z()
	if err := worldcontent.WriteWorld(cfg); err != nil {
		log.Fatalf("Error saving world. %v\n", err)
	}
}

func (s *GamingState) clickCallback(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
	if action == glfw.Press {
		if button == glfw.MouseButtonRight {
			b := s.bottomBar.GetSelectedBlock()
			if b != block.Air {
				if s.world.PlaceBlock(s.camera, b) {
					s.inventory.RemoveBottomBar(s.bottomBar.GetSelectedIndex())
					s.bottomBar.ReBuild()
				}
			}
		} else if button == glfw.MouseButtonLeft {
			 broken, brokenBlock := s.world.BreakBlock(s.camera)
			 if broken {
			 	s.inventory.Add(brokenBlock)
			 	s.bottomBar.ReBuild()
			 }
		}
	}
}

func (s *GamingState) mouseMoveCallback(w *glfw.Window, xpos float64, ypos float64) {
	x, y := s.display.GLPos(xpos, ypos)
	s.player.Entity.Rotation = mgl32.Vec3{0, -x*s.settings.cameraSensitivity, 0}
	s.camera.Rotation = mgl32.Vec3{y*s.settings.cameraSensitivity, s.camera.Rotation.Y(), s.camera.Rotation.Z()}
	pointed, _ := s.world.GetPointedBlock(s.camera)
	s.cubeOutline.Position = mgl32.Vec3{float32(pointed.X), float32(pointed.Y), float32(pointed.Z)}
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
			if s.worldContent.Config().Creative {
				now := time.Now()
				if now.Sub(s.keyPressed.previousSpacePressed) < s.settings.doublePressDelay {
					s.player.InFlight = !s.player.InFlight
					s.player.Speed[1] = 0
				}
				s.keyPressed.previousSpacePressed = now
			}
		} else if key == glfw.KeyLeftShift {
			s.keyPressed.shiftPressed = true
		} else if key == glfw.KeyEscape {
			s.changeState <- state.Switch{ID: state.GameMenu}
		} else if key == glfw.KeyE {
			s.changeState <- state.Switch{ID: state.Inventory}
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
		} else if key == glfw.KeyLeftShift {
			s.keyPressed.shiftPressed = false
		}
	}
}

func (s *GamingState) scrollCallBack(w *glfw.Window, xoff float64, yoff float64) {
	s.scroll += yoff
	offset := int(s.scroll / s.settings.scrollStep)
	if offset != 0 {
		s.scroll -= float64(offset) * s.settings.scrollStep
		s.bottomBar.OffsetSelectedItem(offset)
	}
}

// Render renders all objects on the screen
func (s *GamingState) Render(renderer *render.MasterRenderer) {
	// why here?
	s.camera.LockOnPlayer(s.player)
	// r.ProcessEntity(player.Entity)
	renderer.ProcessEntities(s.world.GetChunks(s.camera))
	renderer.SetCamera(s.camera)
	renderer.ProcessGui(s.cursor)
	renderer.ProcessOutlineModel(s.cubeOutline)
	s.bottomBar.Render(renderer)
}
// NextFrame makes time pass to move to the next frame of the game
func (s *GamingState) NextFrame() {
Loop:
	for {
		select {
		case chunk := <-s.chunkLoader.Chunks():
			s.world.AddChunk(chunk)
		default:
			break Loop
		}
	}
	forward := s.keyPressed.wPressed
	backward := s.keyPressed.sPressed
	right := s.keyPressed.dPressed
	left := s.keyPressed.aPressed
	up := s.keyPressed.spacePressed
	down := s.keyPressed.shiftPressed
	touchGround := s.world.TouchesGround(s.player)
	s.player.Move(forward, backward, up, down, touchGround, right, left)
	s.world.MovePlayer(s.player, forward, backward, up, touchGround)
}
// Update is called every second
func (s *GamingState) Update() {
	s.world.LoadChunks(s.player.Entity.Position)
}

// pause is used when menu is opened
func (s *GamingState) pause() {
	s.keyPressed = keyPressed{}
}

func (s *GamingState) Resize(aspectRatio float32) {
	s.world.Resize(aspectRatio)
	s.bottomBar.Resize(aspectRatio)
}