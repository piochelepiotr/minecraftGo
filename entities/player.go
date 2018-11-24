package entities

import (
	"github.com/go-gl/mathgl/mgl32"
)

type Player struct {
	Entity Entity
}

func (p *Player) MoveForward(dist float32) {
	mat := mgl32.HomogRotate3DY(p.Entity.Rotation.Y())
	forward := mgl32.Vec4{0, 0, -dist, 1}
	forward = mat.Mul4x1(forward)
	p.Entity.Position = p.Entity.Position.Add(forward.Vec3())
}
