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

type Entities []Entity

func (e Entities) Len() int {
	return len(e)
}

func (e Entities) Less(i, j int) bool {
	if e[i].Position[0] < e[j].Position[0] {
		return true
	}
	if e[j].Position[0] < e[i].Position[0] {
		return false
	}
	if e[i].Position[1] < e[j].Position[1] {
		return true
	}
	if e[j].Position[1] < e[i].Position[1] {
		return false
	}
	return e[i].Position[2] < e[j].Position[2]
}

func (e Entities) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

