package game

import (
	"fmt"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/piochelepiotr/minecraftGo/entities"
	"github.com/piochelepiotr/minecraftGo/guis"
	"github.com/piochelepiotr/minecraftGo/loader"
	pmenu "github.com/piochelepiotr/minecraftGo/menu"
	"github.com/piochelepiotr/minecraftGo/models"
	"github.com/piochelepiotr/minecraftGo/render"
	"github.com/piochelepiotr/minecraftGo/state"
	pworld "github.com/piochelepiotr/minecraftGo/world"
	"os"
	"time"
)

type keyPressed struct {
	wPressed bool
	aPressed bool
	dPressed bool
	sPressed bool
	spacePressed bool
}

type GameState struct {
	world *pworld.World
	player *entities.Player
	camera *entities.Camera
	cursor guis.GuiTexture
	light *entities.Light
	menu *pmenu.Menu
	keyPressed keyPressed
	changeState chan state.StateID
	menuOpened bool
	display *render.DisplayManager
	chunkLoader *pworld.ChunkLoader
}

func Run(aspectRatio float32, changeState chan state.StateID, display *render.DisplayManager) {
	generator := pworld.NewGenerator()
	chunkLoader := pworld.NewChunkLoader(generator)
	world := pworld.CreateWorld(generator)
	defer world.Close()
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


	menu := pmenu.CreateMenu(aspectRatio)
	menu.AddItem("Resume game", func() { changeState <- state.Game })
	menu.AddItem("Exit game", func() { os.Exit(0) })
	menu.AddItem("Watch YouTube", func() {})
	menu.AddItem("Go to Website", func() {})

	gameState := &GameState{
		world: world,
		player: player,
		camera: camera,
		cursor: loader.LoadGuiTexture("textures/cursor.png", mgl32.Vec2{0, 0}, mgl32.Vec2{0.02, 0.03}),
		menu: menu,
		light: light,
		changeState: changeState,
		display: display,
		chunkLoader: chunkLoader,
	}


	// set callbacks

	clickCallback := func(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
		if action == glfw.Press {
			if button == glfw.MouseButtonRight {
				gameState.RightClick()
			} else if button == glfw.MouseButtonLeft {
				gameState.LeftClick()
			}
		}
	}

	mouseMoveCallback := func(w *glfw.Window, xpos float64, ypos float64) {
		x, y := display.GLPos(xpos, ypos)
		gameState.MouseMove(x, y)
	}

	keyCallback := func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		if action == glfw.Press {
			gameState.KeyPressed(key)
		} else if action == glfw.Release {
			gameState.KeyReleased(key)
		}
	}

	display.Window.SetKeyCallback(keyCallback)
	display.Window.SetCursorPosCallback(mouseMoveCallback)
	display.Window.SetMouseButtonCallback(clickCallback)
	gameState.run()
}

func (g *GameState) run() {
	renderer := render.CreateMasterRenderer()
	defer renderer.CleanUp()
	defer loader.CleanUp()


	updateTicker := time.NewTicker(time.Second)
	defer updateTicker.Stop()

	frames := 0
	for !g.display.Window.ShouldClose() {
		select {
		case <-updateTicker.C:
			start := time.Now()
			g.Update()
			stopTime := time.Now().Sub(start)
			fmt.Println(stopTime)
			fmt.Printf("FPS is %d\n", frames)
			frames = 0
		case stateID := <- g.changeState:
			if stateID == state.Game {
				g.display.Window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
				g.CloseMenu()
			} else if stateID == state.GameMenu {
				g.display.Window.SetInputMode(glfw.CursorMode, glfw.CursorNormal)
				g.OpenMenu()
			}
		case chunk := <-g.chunkLoader.LoadedChunk:
			g.world.AddChunk(chunk)
		default:
			frames += 1
			g.Render(renderer)
			g.NextFrame()
			g.display.UpdateDisplay()
			glfw.PollEvents()
		}
	}

}

// update is called every second
func (g *GameState) Update() {
	g.world.LoadChunks(g.player.Entity.Position, true)
}

func (g *GameState) NextFrame() {
	forward := g.keyPressed.wPressed
	backward := g.keyPressed.sPressed
	right := g.keyPressed.dPressed
	left := g.keyPressed.aPressed
	jump := g.keyPressed.spacePressed
	touchGround := g.world.TouchesGround(g.player)
	g.player.Move(forward, backward, jump, touchGround, right, left)
	g.world.MovePlayer(g.player, forward, backward, jump, touchGround)
}

func (g *GameState) Render(renderer *render.MasterRenderer) {
	g.camera.LockOnPlayer(g.player)
	// r.ProcessEntity(player.Entity)
	renderer.ProcessEntities(g.world.GetChunks(g.camera))
	renderer.ProcessGui(g.cursor)
	renderer.Render(g.light, g.camera)
	if g.menuOpened {
		g.menu.Render(renderer)
	}
}

func (g *GameState) MouseMove(x, y float32) {
	if g.menuOpened {
		g.menu.ComputeSelectedItem(x, y)
	} else {
		g.player.Entity.Rotation = mgl32.Vec3{0, -x, 0}
		g.camera.Rotation = mgl32.Vec3{y, g.camera.Rotation.Y(), g.camera.Rotation.Z()}
	}
}

func (g *GameState) KeyPressed(key glfw.Key) {
	if key == glfw.KeyW {
		g.keyPressed.wPressed = true
	} else if key == glfw.KeyA {
		g.keyPressed.aPressed = true
	} else if key == glfw.KeyD {
		g.keyPressed.dPressed = true
	} else if key == glfw.KeyS {
		g.keyPressed.sPressed = true
	} else if key == glfw.KeySpace {
		g.keyPressed.spacePressed = true
	} else if key == glfw.KeyEscape {
		if !g.menuOpened {
			g.changeState <- state.GameMenu
		}
	}
}

func (g *GameState) KeyReleased(key glfw.Key) {
	if key == glfw.KeyW {
		g.keyPressed.wPressed = false
	} else if key == glfw.KeyA {
		g.keyPressed.aPressed = false
	} else if key == glfw.KeyD {
		g.keyPressed.dPressed = false
	} else if key == glfw.KeyS {
		g.keyPressed.sPressed = false
	} else if key == glfw.KeySpace {
		g.keyPressed.spacePressed = false
	}
}

func (g *GameState) OpenMenu() {
	g.menuOpened = true
}

func (g *GameState) CloseMenu() {
	g.menuOpened = false
}

func (g *GameState) LeftClick() {
	if g.menuOpened {
		g.menu.LeftClick()
	} else {
		g.world.ClickOnBlock(g.camera, false)
	}
}

func (g *GameState) RightClick() {
	fmt.Println(g.camera.Position)
	if !g.menuOpened {
		g.world.ClickOnBlock(g.camera, true)
	}
}
