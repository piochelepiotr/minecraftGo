package main

import (
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/piochelepiotr/minecraftGo/entities"
	"github.com/piochelepiotr/minecraftGo/loader"
	pmenu "github.com/piochelepiotr/minecraftGo/menu"
	"github.com/piochelepiotr/minecraftGo/models"
	"github.com/piochelepiotr/minecraftGo/render"
	pworld "github.com/piochelepiotr/minecraftGo/world"
	"os"
	"time"
)

const windowWidth = 800
const windowHeight = 600
const aspectRatio = windowWidth / windowHeight

func main() {
	d := render.DisplayManager{WindowWidth: windowWidth, WindowHeight: windowHeight}
	defer d.CloseDisplay()
	d.CreateDisplay()

	model := loader.LoadObjModel("objects/steve.obj")
	defer loader.CleanUp()
	r := render.CreateMasterRenderer()
	defer r.CleanUp()
	t := loader.LoadModelTexture("textures/skin.png")
	cubeTexture := loader.LoadModelTexture("textures/textures2.png")
	cubeTexture.NumberOfRows = 16
	texturedModel := models.TexturedModel{
		ModelTexture: t,
		RawModel:     model,
	}
	world := pworld.CreateWorld(cubeTexture)
	for x := -1; x < 2; x++ {
		for y := 0; y < 2; y++ {
			for z := -1; z < 2; z++ {
				world.LoadChunk(x*pworld.ChunkSize, y*pworld.ChunkSize, z*pworld.ChunkSize)
			}
		}
	}
	camera := entities.CreateCamera(-50, 30, -50, -0.2, 1.8)
	camera.Rotation = mgl32.Vec3{0, 0, 0}

	light := entities.Light{
		Position: mgl32.Vec3{5, 5, 5},
		Colour:   mgl32.Vec3{1, 1, 1},
	}

	entity := entities.Entity{
		Position:      mgl32.Vec3{0, float32(world.GetHeight(int(0), int(0))), 0},
		Rotation:      mgl32.Vec3{0, 0, 0},
		TexturedModel: texturedModel,
	}
	player := entities.Player{
		Entity: entity,
	}

	cursor := loader.LoadGuiTexture("textures/cursor.png", mgl32.Vec2{0, 0}, mgl32.Vec2{0.02, 0.03})

	menu := pmenu.CreateMenu(aspectRatio)
	menu.Opened = false
	menu.AddItem("Resume game")
	menu.AddItem("Exit game")
	menu.AddItem("Watch YouTube")
	menu.AddItem("Go to Website")

	movePlayer := func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		if !menu.Opened {
			if key == glfw.KeyEscape {
				menu.Opened = true
				d.Window.SetInputMode(glfw.CursorMode, glfw.CursorNormal)
			} else if key == glfw.KeyW {
				//player.Accelerate(4)
			} else if key == glfw.KeyS {
				//player.Accelerate(-4)
			} else if key == glfw.KeyD {
				//player.Entity.IncreasePosition(0, -0.1, 0)
			} else if key == glfw.KeyU {
				//player.Entity.IncreasePosition(0, 0.1, 0)
			} else if key == glfw.KeySpace {
				//player.Jump()
			}

			//else if key == glfw.KeyA {
			//	player.Entity.IncreaseRotation(0.0, -0.1, 0.0)
			//} else if key == glfw.KeyW {
			//	player.Entity.IncreasePosition(0.0, 0.1, 0.0)
			//} else if key == glfw.KeyS {
			//	player.Entity.IncreasePosition(0.0, -0.1, 0.0)
			//} else if key == glfw.KeyJ {
			//	player.Entity.IncreasePosition(0.1, 0.0, 0.0)
			//} else if key == glfw.KeyL {
			//	player.Entity.IncreasePosition(-0.1, 0.0, 0.0)
			//} else if key == glfw.KeyI {
			//	player.Entity.IncreasePosition(0.0, 0.0, 0.1)
			//} else if key == glfw.KeyK {
			//	player.Entity.IncreasePosition(0.0, 0.0, -0.1)
			//}
		}
	}

	menuSelectItem := func(w *glfw.Window, xpos float64, ypos float64) {
		x, y := d.GLPos(xpos, ypos)
		if menu.Opened {
			menu.ComputeSelectedItem(x, y)
		} else {
			player.Entity.Rotation = mgl32.Vec3{0, -x, 0}
			camera.Rotation = mgl32.Vec3{y, camera.Rotation.Y(), camera.Rotation.Z()}
		}
	}

	menuClick := func(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
		if menu.Opened {
			if button == glfw.MouseButtonLeft {
				if menu.SelectedItem == 0 {
					menu.Opened = false
					d.Window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
				} else if menu.SelectedItem == 1 {
					os.Exit(0)
				}
			}
		} else {
			if button == glfw.MouseButtonLeft && action == glfw.Press {
				world.ClickOnBlock(&camera)
			}
		}
	}

	resizeWindow := func(w *glfw.Window, width int, height int) {
		d.Resize(width, height)
	}

	d.Window.SetKeyCallback(movePlayer)
	d.Window.SetCursorPosCallback(menuSelectItem)
	d.Window.SetMouseButtonCallback(menuClick)
	d.Window.SetSizeCallback(resizeWindow)

	world.LoadChunks(player.Entity.Position)

	updateTicker := time.NewTicker(time.Second)
	defer updateTicker.Stop()

	for !d.Window.ShouldClose() {
		select {
			case <-updateTicker.C:
				world.LoadChunks(player.Entity.Position)
			default:
				//world.LoadChunks(player.Entity.Position)
				camera.LockOnPlayer(player)
				// r.ProcessEntity(player.Entity)
				r.ProcessEntities(world.GetChunks())
				r.ProcessGui(cursor)
				r.ProcessMenu(menu)
				r.Render(light, camera)
				d.UpdateDisplay()
				forward := d.Window.GetKey(glfw.KeyW) == glfw.Press
				backward := d.Window.GetKey(glfw.KeyS) == glfw.Press
				jump := d.Window.GetKey(glfw.KeySpace) == glfw.Press
				touchGround := world.TouchesGround(&player)
				//player.Entity.IncreaseRotation(0.0, 0.1, 0.0)
				//player.Move(d.Window.GetKey(glfw.KeyW) == glfw.Press, move)
				player.Move(forward, backward, jump, touchGround)
				world.MovePlayer(&player, forward, backward, jump, touchGround)
				glfw.PollEvents()
		}
	}
}
