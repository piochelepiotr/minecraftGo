package entities

import (
	"github.com/go-gl/mathgl/mgl32"
)


type Camera struct {
    Entity Entity
    FollowDistance float32
    Height float32
}

func CreateCamera(x, y, z float32) Camera {
    e := Entity{
        Position:mgl32.Vec3{x, y, z},
        Rotation:mgl32.Vec3{0, 0, 0},
    }
    return Camera{Entity:e}
}

func (c *Camera) lockOnPlayer(player Player) {
    c.Height = 1.8
    c.FollowDistance = -5
   m := mgl32.Ident4()
   rotY := mgl32.Rotate3DY(player.Entity.Rotation.Y())
   cameraShift := mgl32.Vec3{0, 0, c.FollowDistance}
   movement := rotY.Mul3x1(cameraShift)
   c.Position = mgl32.Vec3{
       player.Entity.X() + movement.X(),
       player.Entity.Y() + c.Height,
       player.Entity.Z() + movement.Z(),
   }
   c.Entity.Rotation.Elem(c.Entity.Rotation.X(), -player.Entity.Rotation.Y(), c.Entity.Rotation.Z())
}
