package main

import (
    "github.com/piochelepiotr/minecraftGo/shaders"
    "github.com/piochelepiotr/minecraftGo/renderEngine"
    "github.com/piochelepiotr/minecraftGo/textures"
    "github.com/piochelepiotr/minecraftGo/models"
)

const windowWidth = 800
const windowHeight = 600

func main() {
    var vertices = []float32{
        -0.5, 0.5, 0,//top left
        -0.5, -0.5, 0,//bottom left
        0.5, -0.5, 0,//bottom right
        0.5, 0.5, 0,//top right
    }
    var indices = []uint32{
        0, 1, 3,
        3, 1, 2,
    }
    var textureCoord = []float32{
        0,0,
        0,1,
        1,1,
        1,0,
    }

    d := renderEngine.DisplayManager{WindowWidth:windowWidth, WindowHeight:windowHeight}
    defer d.CloseDisplay()
    d.CreateDisplay()

    model := renderEngine.LoadToVAO(vertices, textureCoord, indices)
    s := shaders.CreateStaticShader()
    textureID, err := renderEngine.LoadTexture("res/pic.png")
    if err != nil {
        panic(err)
    }
    t := textures.ModelTexture{Id:textureID}
    texturedModel := models.TexturedModel{ModelTexture:t, RawModel:model}

	for !d.Window.ShouldClose() {
        renderEngine.Prepare()
        s.Start()
        renderEngine.Render(texturedModel)
        s.Stop()
        d.UpdateDisplay()
	}
    s.CleanUp()
    renderEngine.CleanUp()
}


