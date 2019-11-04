package models
import (
    "github.com/piochelepiotr/minecraftGo/textures"
)

type TexturedModel struct {
    RawModel RawModel
    ModelTexture textures.ModelTexture
    Transparent bool
}
