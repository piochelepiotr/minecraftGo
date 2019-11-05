package game

import (
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
	renderer *render.MasterRenderer
	cursor guis.GuiTexture
	light *entities.Light
	menu *pmenu.Menu
	keyPressed keyPressed
	changeState chan<- state.StateID
}

func NewGameState(aspectRatio float32, changeState chan<- state.StateID) *GameState {
	renderer := render.CreateMasterRenderer()

	world := pworld.CreateWorld()
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
	world.LoadChunks(player.Entity.Position)
	world.PlacePlayerOnGround(player)


	menu := pmenu.CreateMenu(aspectRatio)
	menu.Opened = false
	menu.AddItem("Resume game")
	menu.AddItem("Exit game")
	menu.AddItem("Watch YouTube")
	menu.AddItem("Go to Website")


	return &GameState{
		world: world,
		player: player,
		camera: camera,
		renderer: renderer,
		cursor: loader.LoadGuiTexture("textures/cursor.png", mgl32.Vec2{0, 0}, mgl32.Vec2{0.02, 0.03}),
		menu: menu,
		light: light,
		changeState: changeState,
	}
}

func (g *GameState) Close() {
	loader.CleanUp()
	g.renderer.CleanUp()
}

// update is called every second
func (g *GameState) Update() {
	g.world.LoadChunks(g.player.Entity.Position)
}

func (g *GameState) NextFrame() {
	g.camera.LockOnPlayer(g.player)
	// r.ProcessEntity(player.Entity)
	g.renderer.ProcessEntities(g.world.GetChunks())
	g.renderer.ProcessGui(g.cursor)
	g.renderer.ProcessMenu(g.menu)
	g.renderer.Render(g.light, g.camera)
	forward := g.keyPressed.wPressed
	backward := g.keyPressed.sPressed
	right := g.keyPressed.dPressed
	left := g.keyPressed.aPressed
	jump := g.keyPressed.spacePressed
	touchGround := g.world.TouchesGround(g.player)
	g.player.Move(forward, backward, jump, touchGround, right, left)
	g.world.MovePlayer(g.player, forward, backward, jump, touchGround)
}

func (g *GameState) MouseMove(x, y float32) {
	if g.menu.Opened {
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
		if !g.menu.Opened {
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
	g.menu.Opened = true
}

func (g *GameState) CloseMenu() {
	g.menu.Opened = false
}

func (g *GameState) LeftClick() {
	if g.menu.Opened {
		if g.menu.SelectedItem == 0 {
			g.changeState <- state.Game
		} else if g.menu.SelectedItem == 1 {
			os.Exit(0)
		}
	} else {
		g.world.ClickOnBlock(g.camera, false)
	}
}

func (g *GameState) RightClick() {
	if !g.menu.Opened {
		g.world.ClickOnBlock(g.camera, true)
	}
}
