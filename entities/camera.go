package entities

import (
	//"github.com/go-gl/mathgl/mgl32"
)


type Camera struct {
    Entity Entity
    FollowDistance float32
    Height float32
}

func (c *Camera) lockOnPlayer(player Player) {
   // m := mgl32.Ident4()
}
