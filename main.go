package main

import (
    "github.com/piochelepiotr/minecraftGo/shaders"
    "github.com/piochelepiotr/minecraftGo/renderEngine"
    "github.com/piochelepiotr/minecraftGo/textures"
    "github.com/piochelepiotr/minecraftGo/models"
    "github.com/piochelepiotr/minecraftGo/entities"
	"github.com/go-gl/mathgl/mgl32"
)

const windowWidth = 800
const windowHeight = 600

func main() {
    d := renderEngine.DisplayManager{WindowWidth:windowWidth, WindowHeight:windowHeight}
    defer d.CloseDisplay()
    d.CreateDisplay()

    model := renderEngine.LoadObjModel("objects/steve.obj")
    s := shaders.CreateStaticShader()
    textureID, err := renderEngine.LoadTexture("textures/skin.png")
    if err != nil {
        panic(err)
    }
    t := textures.ModelTexture{
        Id:textureID,
        Reflectivity: 1,
        ShineDamper: 10,
    }
    texturedModel := models.TexturedModel{ModelTexture:t, RawModel:model}
    entity := entities.Entity{
        Position: mgl32.Vec3{0, -1.5, -5},
        Rotation: mgl32.Vec3{0, 0, 0},
    }
    player := entities.Player{
        Entity: entity,
        TexturedModel: texturedModel,
    }
    camera := entities.CreateCamera(0, 0, 1)

    light := entities.Light{
        Position: mgl32.Vec3{5, 5, 5},
        Colour: mgl32.Vec3{1, 1, 1},
    }


	for !d.Window.ShouldClose() {
        camera.LockOnPlayer(player)
        renderEngine.Prepare()
        s.Program.Start()
        s.LoadLight(light)
        renderEngine.Render(player, camera, s)
        s.Program.Stop()
        d.UpdateDisplay()
        player.Entity.IncreaseRotation(0.0, 0.1, 0.0)
	}
    s.Program.CleanUp()
    renderEngine.CleanUp()
}


