package main

import (
    "github.com/piochelepiotr/minecraftGo/renderEngine"
    "github.com/piochelepiotr/minecraftGo/models"
    "github.com/piochelepiotr/minecraftGo/entities"
    pworld "github.com/piochelepiotr/minecraftGo/world"
	"github.com/go-gl/mathgl/mgl32"
    "github.com/go-gl/glfw/v3.2/glfw"
)

const windowWidth = 800
const windowHeight = 600

func main() {
    d := renderEngine.DisplayManager{WindowWidth:windowWidth, WindowHeight:windowHeight}
    defer d.CloseDisplay()
    d.CreateDisplay()

    model := renderEngine.LoadObjModel("objects/steve.obj")
    defer renderEngine.CleanUp()
    r := renderEngine.CreateMasterRenderer()
    defer r.CleanUp()
    t := renderEngine.LoadTexture("textures/skin.png")
    cubeTexture := renderEngine.LoadTexture("textures/textures.png")
    cubeTexture.NumberOfRows = 2
    texturedModel := models.TexturedModel{
        ModelTexture:t,
        RawModel:model,
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
        Colour: mgl32.Vec3{1, 1, 1},
    }

    camera.IncreaseRotation(0, 2, 0)

    entity := entities.Entity{
        Position: mgl32.Vec3{0, float32(world.GetHeight(int(0), int(0))), 0},
        Rotation: mgl32.Vec3{0, 0, 0},
        TexturedModel: texturedModel,
    }
    player := entities.Player{
        Entity: entity,
    }

    movePlayer := func (w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
        if (key == glfw.KeyD) {
            player.Entity.IncreaseRotation(0.0, 0.1, 0.0)
        } else if (key == glfw.KeyA) {
            player.Entity.IncreaseRotation(0.0, -0.1, 0.0)
        } else if (key == glfw.KeyW) {
            player.Entity.IncreasePosition(0.0, 0.1, 0.0)
        } else if (key == glfw.KeyS) {
            player.Entity.IncreasePosition(0.0, -0.1, 0.0)
        } else if (key == glfw.KeyJ) {
            player.Entity.IncreasePosition(0.1, 0.0, 0.0)
        } else if (key == glfw.KeyL) {
            player.Entity.IncreasePosition(-0.1, 0.0, 0.0)
        } else if (key == glfw.KeyI) {
            player.Entity.IncreasePosition(0.0, 0.0, 0.1)
        } else if (key == glfw.KeyK) {
            player.Entity.IncreasePosition(0.0, 0.0, -0.1)
        }
    }

    d.Window.SetKeyCallback(movePlayer)


	for !d.Window.ShouldClose() {
        camera.LockOnPlayer(player)
        r.ProcessEntity(player.Entity)
        r.ProcessEntities(world.GetChunks())
        r.Render(light, camera)
        d.UpdateDisplay()
        //player.Entity.IncreaseRotation(0.0, 0.1, 0.0)
        glfw.PollEvents()
	}
}


