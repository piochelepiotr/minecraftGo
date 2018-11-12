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
   // m := mgl32.Ident4()
}
