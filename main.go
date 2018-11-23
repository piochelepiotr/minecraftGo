package main

import (
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/piochelepiotr/minecraftGo/entities"
	"github.com/piochelepiotr/minecraftGo/fontRendering"
	pguis "github.com/piochelepiotr/minecraftGo/guis"
	pmenu "github.com/piochelepiotr/minecraftGo/menu"
	"github.com/piochelepiotr/minecraftGo/models"
	"github.com/piochelepiotr/minecraftGo/renderEngine"
	pworld "github.com/piochelepiotr/minecraftGo/world"
)

const windowWidth = 800
const windowHeight = 600
const aspectRatio = windowWidth / windowHeight

func main() {
	d := renderEngine.DisplayManager{WindowWidth: windowWidth, WindowHeight: windowHeight}
	defer d.CloseDisplay()
	d.CreateDisplay()

	model := renderEngine.LoadObjModel("objects/steve.obj")
	defer renderEngine.CleanUp()
	r := renderEngine.CreateMasterRenderer()
	defer r.CleanUp()
	fontRenderer := fontRendering.CreateFontRenderer()
	defer fontRenderer.CleanUp()
	t := renderEngine.LoadModelTexture("textures/skin.png")
	cubeTexture := renderEngine.LoadModelTexture("textures/textures.png")
	cubeTexture.NumberOfRows = 2
	texturedModel := models.TexturedModel{
		ModelTexture: t,
		RawModel:     model,
	}
	world := pworld.CreateWorld(cubeTexture)
	world.LoadChunk(0, 0, 0)
	world.LoadChunk(-pworld.ChunkSize, 0, 0)
	world.LoadChunk(0, 0, -pworld.ChunkSize)
	world.LoadChunk(-pworld.ChunkSize, 0, -pworld.ChunkSize)
	camera := entities.CreateCamera(-50, 30, -50, -0.5, 1.8)
	camera.Rotation = mgl32.Vec3{0, 0, 0}

	light := entities.Light{
		Position: mgl32.Vec3{5, 5, 5},
		Colour:   mgl32.Vec3{1, 1, 1},
	}

	camera.IncreaseRotation(0, 2, 0)

	entity := entities.Entity{
		Position:      mgl32.Vec3{0, float32(world.GetHeight(int(0), int(0))), 0},
		Rotation:      mgl32.Vec3{0, 0, 0},
		TexturedModel: texturedModel,
	}
	player := entities.Player{
		Entity: entity,
	}

	guis := make([]pguis.GuiTexture, 0)
	cursor := renderEngine.LoadGuiTexture("textures/cursor.png", mgl32.Vec2{0, 0}, mgl32.Vec2{0.02, 0.03})
	guis = append(guis, cursor)

	menu := pmenu.CreateMenu(aspectRatio)
	menu.Opened = true
	menu.AddItem("Resume game")
	menu.AddItem("Exit game")
	menu.AddItem("Watch YouTube")
	menu.AddItem("Go to Website")
	guis = append(guis, menu.GetMenuItems()...)
	fontRenderer.LoadTexts(menu.GetMenuTexts())

	guiRenderer := renderEngine.CreateGuiRenderer()
	defer guiRenderer.CleanUp()

	movePlayer := func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		if !menu.Opened {
			if key == glfw.KeyD {
				player.Entity.IncreaseRotation(0.0, 0.1, 0.0)
			} else if key == glfw.KeyA {
				player.Entity.IncreaseRotation(0.0, -0.1, 0.0)
			} else if key == glfw.KeyW {
				player.Entity.IncreasePosition(0.0, 0.1, 0.0)
			} else if key == glfw.KeyS {
				player.Entity.IncreasePosition(0.0, -0.1, 0.0)
			} else if key == glfw.KeyJ {
				player.Entity.IncreasePosition(0.1, 0.0, 0.0)
			} else if key == glfw.KeyL {
				player.Entity.IncreasePosition(-0.1, 0.0, 0.0)
			} else if key == glfw.KeyI {
				player.Entity.IncreasePosition(0.0, 0.0, 0.1)
			} else if key == glfw.KeyK {
				player.Entity.IncreasePosition(0.0, 0.0, -0.1)
			}
		}
	}

	menuSelectItem := func(w *glfw.Window, xpos float64, ypos float64) {
		xpos = xpos / float64(d.WindowWidth)
		ypos = ypos / float64(d.WindowHeight)
		xpos = xpos - 0.5
		ypos = ypos - 0.5
		menu.ComputeSelectedItem(xpos, ypos)
	}

	d.Window.SetKeyCallback(movePlayer)
	d.Window.SetCursorPosCallback(menuSelectItem)

	for !d.Window.ShouldClose() {
		camera.LockOnPlayer(player)
		r.ProcessEntity(player.Entity)
		r.ProcessEntities(world.GetChunks())
		guis = append(guis[:1], menu.GetMenuItems()...)
		r.Render(light, camera)
		guiRenderer.Render(guis)
		fontRenderer.Render()
		d.UpdateDisplay()
		//player.Entity.IncreaseRotation(0.0, 0.1, 0.0)
		glfw.PollEvents()
	}
}
