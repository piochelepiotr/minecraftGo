package entities
import (
    "github.com/piochelepiotr/minecraftGo/models"
)

type Player struct {
    Entity Entity
    TexturedModel models.TexturedModel
    TextureIndex uint32
}

