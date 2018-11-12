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
    //var vertices = []float32{
    //    -0.5, 0.5, 0,//top left
    //    -0.5, -0.5, 0,//bottom left
    //    0.5, -0.5, 0,//bottom right
    //    0.5, 0.5, 0,//top right
    //}
    //var indices = []uint32{
    //    0, 1, 3,
    //    3, 1, 2,
    //}
    //var textureCoord = []float32{
    //    0,0,
    //    0,1,
    //    1,1,
    //    1,0,
    //}
    //var vertices = []float32{
    //    -0.5,0.5,-0.5,
    //    -0.5,-0.5,-0.5,
    //    0.5,-0.5,-0.5,
    //    0.5,0.5,-0.5,

    //    -0.5,0.5,0.5,
    //    -0.5,-0.5,0.5,
    //    0.5,-0.5,0.5,
    //    0.5,0.5,0.5,

    //    0.5,0.5,-0.5,
    //    0.5,-0.5,-0.5,
    //    0.5,-0.5,0.5,
    //    0.5,0.5,0.5,

    //    -0.5,0.5,-0.5,
    //    -0.5,-0.5,-0.5,
    //    -0.5,-0.5,0.5,
    //    -0.5,0.5,0.5,

    //    -0.5,0.5,0.5,
    //    -0.5,0.5,-0.5,
    //    0.5,0.5,-0.5,
    //    0.5,0.5,0.5,

    //    -0.5,-0.5,0.5,
    //    -0.5,-0.5,-0.5,
    //    0.5,-0.5,-0.5,
    //    0.5,-0.5,0.5,
    //}
    //var indices = []uint32{
    //    0,1,3,
    //    3,1,2,
    //    4,5,7,
    //    7,5,6,
    //    8,9,11,
    //    11,9,10,
    //    12,13,15,
    //    15,13,14,
    //    16,17,19,
    //    19,17,18,
    //    20,21,23,
    //    23,21,22,
    //}
    //var textureCoord = []float32{
    //    0,0,
    //    0,1,
    //    1,1,
    //    1,0,
    //    0,0,
    //    0,1,
    //    1,1,
    //    1,0,
    //    0,0,
    //    0,1,
    //    1,1,
    //    1,0,
    //    0,0,
    //    0,1,
    //    1,1,
    //    1,0,
    //    0,0,
    //    0,1,
    //    1,1,
    //    1,0,
    //    0,0,
    //    0,1,
    //    1,1,
    //    1,0,
    //}

    d := renderEngine.DisplayManager{WindowWidth:windowWidth, WindowHeight:windowHeight}
    defer d.CloseDisplay()
    d.CreateDisplay()

    model := renderEngine.LoadObjModel("objects/steve.obj") //renderEngine.LoadToVAO(vertices, textureCoord, indices)
    s := shaders.CreateStaticShader()
    textureID, err := renderEngine.LoadTexture("textures/skin.png")
    if err != nil {
        panic(err)
    }
    t := textures.ModelTexture{Id:textureID}
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

    camera.LockOnPlayer(player)

	for !d.Window.ShouldClose() {
        renderEngine.Prepare()
        s.Program.Start()
        renderEngine.Render(player, camera, s)
        s.Program.Stop()
        d.UpdateDisplay()
        player.Entity.IncreaseRotation(0.0, 0.01, 0.0)
	}
    s.Program.CleanUp()
    renderEngine.CleanUp()
}


