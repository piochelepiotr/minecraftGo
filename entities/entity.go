package entities

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/piochelepiotr/minecraftGo/game_engine/models"
)

type Entity struct {
	Position      mgl32.Vec3
	Rotation      mgl32.Vec3
	TexturedModel models.TexturedModel
	TextureIndex  uint32
}

func (e *Entity) IncreaseRotation(dx, dy, dz float32) {
	e.Rotation = e.Rotation.Add(mgl32.Vec3{dx, dy, dz})
}

func (e *Entity) IncreasePosition(dx, dy, dz float32) {
	e.Position = e.Position.Add(mgl32.Vec3{dx, dy, dz})
}
