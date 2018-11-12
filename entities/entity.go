package entities

import (
	"github.com/go-gl/mathgl/mgl32"
)


type Entity struct {
    Position mgl32.Vec3
    Rotation mgl32.Vec3
}

func (e *Entity) IncreaseRotation(dx, dy, dz float32) {
    e.Rotation = e.Rotation.Add(mgl32.Vec3{dx, dy, dz})
}

